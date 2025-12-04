package util

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"regexp"
	"strings"
)

// stripMarkdown removes Markdown bold and inline code markup from a string.
func stripMarkdown(input string) string {
	// Strip Markdown bold (**...**)
	boldRegex := regexp.MustCompile(`\*\*(.*?)\*\*`)
	// Strip Markdown inline code (`...`)
	backtickRegex := regexp.MustCompile("`(.*?)`")
	// Strip HTML inline code (<code>...</code>)
	codeRegex := regexp.MustCompile(`<code>(.*?)</code>`)
	// Strip '- ' after a newline
	newlineDashRegex := regexp.MustCompile(`(?m)^\- `)

	result := boldRegex.ReplaceAllString(input, "$1")
	result = backtickRegex.ReplaceAllString(result, "$1")
	result = codeRegex.ReplaceAllString(result, "$1")
	result = newlineDashRegex.ReplaceAllString(result, "")

	return result
}

func printFlags(flags *pflag.FlagSet) {

	m := 0
	flags.VisitAll(func(flag *pflag.Flag) {
		n := 0
		if len(flag.Shorthand) > 0 {
			n += len("  -" + flag.Shorthand + ", --" + flag.Name)
		} else {
			n += len("  --" + flag.Name)
		}
		if flag.Value.Type() != "bool" {
			n += len(" <" + flag.Value.Type() + ">")
		}
		if n > m {
			m = n
		}
	})

	var line strings.Builder

	flags.VisitAll(func(flag *pflag.Flag) {
		line.Reset()
		if len(flag.Shorthand) > 0 {
			line.WriteString(fmt.Sprintf("  -%s, --%s", flag.Shorthand, flag.Name))
		} else {
			line.WriteString(fmt.Sprintf("  --%s", flag.Name))
		}
		if flag.Value.Type() != "bool" {
			line.WriteString(fmt.Sprintf(" <%s>", flag.Value.Type()))
		}
		usage := stripMarkdown(flag.Usage)
		if 2+m+4+len(usage) < 80 {
			line.WriteString(fmt.Sprintf("%s  %s\n", strings.Repeat(" ", m+1-len(line.String())), usage))
		} else {
			line.WriteString(fmt.Sprintf("\n%s\n\n", IndentMultilineString(usage, 6)))
		}
		fmt.Print(line.String())
	})
	fmt.Println()
}

func IndentMultilineString(s string, indent int) string {
	lines := strings.Split(s, "\n")
	indentStr := strings.Repeat(" ", indent)
	for i, line := range lines {
		lines[i] = indentStr + line
	}
	return strings.Join(lines, "\n")
}

func CustomHelpFunc(cmd *cobra.Command, args []string) {
	fmt.Printf("Usage:\n  %s\n\n", cmd.UseLine())
	fmt.Println(stripMarkdown(cmd.Long))

	if cmd.Aliases != nil {
		fmt.Printf("Aliases:\n  %s\n\n", strings.Join(cmd.Aliases, ", "))
	}

	flags := cmd.NonInheritedFlags()
	if flags.HasAvailableFlags() {
		fmt.Println("Flags:")
		printFlags(flags)
	}
	parentFlags := cmd.InheritedFlags()
	if parentFlags.HasAvailableFlags() {
		fmt.Println("Global Flags:")
		printFlags(parentFlags)
	}

	if cmd.HasAvailableSubCommands() {
		fmt.Println("Available Commands:")
		maxCmdName := 0
		for _, n := range cmd.Commands() {
			if len(n.Name()) > maxCmdName {
				maxCmdName = len(n.Name())
			}
		}

		for _, c := range cmd.Commands() {
			fmt.Printf("  %s%s%s\n", c.Name(),
				strings.Repeat(" ", maxCmdName+2-len(c.Name())), c.Short)
		}
	}

	if cmd.HasExample() {
		fmt.Println("Examples:")
		fmt.Println(IndentMultilineString(cmd.Example, 2))
	}
	fmt.Println("\nUse todoister help [command] for more information about a command.")
}
