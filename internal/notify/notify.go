package notify

import (
	"github.com/fatih/color"

	"os"
)

const (
	ErrorReadingInput        = "reading input: "
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
		primary = color.New(color.BgRed, color.Bold)
	case TypeWarning:
		prefix = "⚠️ Warning: "
		primary = color.New(color.BgYellow, color.Bold)
	case TypeInfo:
		prefix = "✨ Info: "
		primary = color.New(color.BgHiYellow, color.Bold)
	case TypeSuccess:
		prefix = "✅ Success: "
		primary = color.New(color.BgGreen, color.Bold)
	case TypeWrite:
		prefix = "✍️ Write: "
		primary = color.New(color.BgCyan, color.Bold)
	case TypeBlock:
		prefix = "🚫 Block: "
		primary = color.New(color.FgHiRed, color.Bold)
	default:
		prefix = "❓ Unknown: "
		primary = color.New(color.FgRed, color.Bold)
	}

	if err != nil {
		primary.Printf("%s%s\n", prefix, message)
		primary.Printf("   Details: %v\n", err)
		os.Exit(1)
	} else {
		primary.Printf("%s%s\n", prefix, message)
	}
}
