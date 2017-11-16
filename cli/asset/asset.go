package asset

import (
	"fmt"
	"math/rand"
	"os"

	. "DNA/cli/common"
	. "DNA/common"
	"DNA/net/httpjsonrpc"

	"github.com/urfave/cli"
)

const (
	RANDBYTELEN = 4
)

func parseAssetName(c *cli.Context) string {
	name := c.String("name")
	if name == "" {
		rbuf := make([]byte, RANDBYTELEN)
		rand.Read(rbuf)
		name = "DNA-" + BytesToHexString(rbuf)
	}

	return name
}

func parseAssetID(c *cli.Context) string {
	asset := c.String("asset")
	if asset == "" {
		fmt.Println("missing flag [--asset]")
		os.Exit(1)
	}

	return asset
}

func parseAddress(c *cli.Context) string {
	if address := c.String("to"); address != "" {
		_, err := ToScriptHash(address)
		if err != nil {
			fmt.Println("invalid receiver address")
			os.Exit(1)
		}
		return address
	} else {
		fmt.Println("missing flag [--to]")
		os.Exit(1)
	}

	return ""
}

func parseHeight(c *cli.Context) int64 {
	height := c.Int64("height")
	if height != -1 {
		return height
	} else {
		fmt.Println("invalid parameter [--height]")
		os.Exit(1)
	}

	return 0
}

func assetAction(c *cli.Context) error {
	if c.NumFlags() == 0 {
		cli.ShowSubcommandHelp(c)
		return nil
	}
	value := c.String("value")
	if value == "" {
		fmt.Println("asset amount is required with [--value]")
		return nil
	}
	var resp []byte
	var err error
	switch {
	case c.Bool("reg"):
		name := parseAssetName(c)
		resp, err = httpjsonrpc.Call(Address(), "registerasset", 0, []interface{}{name, value})
	case c.Bool("issue"):
		assetID := parseAssetID(c)
		address := parseAddress(c)
		resp, err = httpjsonrpc.Call(Address(), "issueasset", 0, []interface{}{assetID, address, value})
	case c.Bool("lock"):
		assetID := parseAssetID(c)
		height := parseHeight(c)
		resp, err = httpjsonrpc.Call(Address(), "lockasset", 0, []interface{}{assetID, value, height})
	case c.Bool("transfer"):
		assetID := parseAssetID(c)
		address := parseAddress(c)
		resp, err = httpjsonrpc.Call(Address(), "sendtoaddress", 0, []interface{}{assetID, address, value})
	default:
		cli.ShowSubcommandHelp(c)
		return nil
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	FormatOutput(resp)

	return nil
}

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:        "asset",
		Usage:       "asset registration, issuance and transfer",
		Description: "With nodectl asset, you could control assert through transaction.",
		ArgsUsage:   "[args]",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "reg, r",
				Usage: "regist a new kind of asset",
			},
			cli.BoolFlag{
				Name:  "issue, i",
				Usage: "issue asset that has been registered",
			},
			cli.BoolFlag{
				Name:  "transfer, t",
				Usage: "transfer asset",
			},
			cli.BoolFlag{
				Name:  "lock",
				Usage: "lock asset",
			},
			cli.StringFlag{
				Name:  "asset, a",
				Usage: "uniq id for asset",
			},
			cli.StringFlag{
				Name:  "name",
				Usage: "asset name",
			},
			cli.StringFlag{
				Name:  "to",
				Usage: "asset to whom",
			},
			cli.StringFlag{
				Name:  "value, v",
				Usage: "asset amount",
				Value: "",
			},
			cli.Int64Flag{
				Name:  "height",
				Usage: "asset lock height",
				Value: -1,
			},
		},
		Action: assetAction,
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			PrintError(c, err, "asset")
			return cli.NewExitError("", 1)
		},
	}
}
