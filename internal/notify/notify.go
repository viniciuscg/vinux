package notify

import (
	"github.com/fatih/color"

	"os"
)

const (
	ErrorInvalidYesOrNoInput = "Invalid input. Please enter 'y' for yes or 'n' for no."
)

type MessageType int

const (
	TypeError MessageType = iota

	TypeWarning

	TypeInfo

	TypeSuccess

	TypeWrite

	TypeBlock
)

func Print(
	msgType MessageType,
	message string,
	err error,
) {
	var prefix string
	var primary *color.Color
	switch msgType {
	case TypeError:
		prefix = "❌ Error: "
		primary = color.New(color.FgHiRed, color.Bold)

	case TypeWarning:
		prefix = "⚠️ Warning: "
		primary = color.New(color.FgHiYellow, color.Bold)

	case TypeInfo:
		prefix = "✨ Info: "
		primary = color.New(color.FgHiCyan, color.Bold)

	case TypeSuccess:
		prefix = "✅ Success: "
		primary = color.New(color.FgHiGreen, color.Bold)

	case TypeWrite:
		prefix = "✍️ Write: "
		primary = color.New(color.FgHiBlue, color.Bold)

	case TypeBlock:
		prefix = "🚫 Block: "
		primary = color.New(color.FgMagenta, color.Bold)

	default:
		prefix = "❓ Unknown: "
		primary = color.New(color.FgWhite, color.Bold)
	}

	if err != nil {
		primary.Printf("%s%s\n", prefix, message)
		primary.Printf("   Details: %v\n", err)
		os.Exit(1)
	} else {
		primary.Printf("%s%s\n", prefix, message)
	}
}
