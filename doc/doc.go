// Copyright 2013-2023 The Cobra Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/layfellow/todoister/cmd"
)

const markdownExtension = ".md"

func replaceAngleBrackets(s string) string {
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "&lt;", "<code>&lt;")
	s = strings.ReplaceAll(s, "&gt;", "&gt;</code>")
	return s
}

func printFlags(buf *bytes.Buffer, flags *pflag.FlagSet) {
	buf.WriteString("<dl>\n")

	flags.VisitAll(func(flag *pflag.Flag) {

		buf.WriteString("  <dt>")
		if len(flag.Shorthand) > 0 {
			buf.WriteString(fmt.Sprintf("<code>-%s</code>, <code>--%s</code>", flag.Shorthand, flag.Name))
		} else {
			buf.WriteString(fmt.Sprintf("<code>--%s</code>", flag.Name))
		}
		if flag.Value.Type() != "bool" {
			buf.WriteString(fmt.Sprintf(" <code>&lt;%s&gt;</code>", flag.Value.Type()))
		}
		buf.WriteString("</dt>\n")
		buf.WriteString("  <dd>" + replaceAngleBrackets(flag.Usage) + "</dd>\n")
	})

	buf.WriteString("</dl>\n\n")
}

func printOptions(buf *bytes.Buffer, cmd *cobra.Command) error {
	// No need to print flags for the root command
	if cmd.Name() == "todoister" {
		return nil
	}

	flags := cmd.NonInheritedFlags()
	if flags.HasAvailableFlags() {
		buf.WriteString("### Flags:\n\n")
		printFlags(buf, flags)
	}

	parentFlags := cmd.InheritedFlags()
	if parentFlags.HasAvailableFlags() {
		buf.WriteString("### Global Flags:\n\n")
		printFlags(buf, parentFlags)
	}
	return nil
}

// CustomGenMarkdown generates markdown output for a command with custom formatting.
func CustomGenMarkdown(cmd *cobra.Command, w io.Writer) error {

	buf := new(bytes.Buffer)
	name := cmd.CommandPath()
	buf.WriteString("## " + name + "\n\n")

	if cmd.Runnable() {
		flagUsage := ""
		if !strings.Contains(cmd.UseLine(), "[flags]") {
			flagUsage = " [flags]"
		}
		buf.WriteString(fmt.Sprintf("```sh\n%s%s\n```\n\n", cmd.UseLine(), flagUsage))
	}
	if len(cmd.Long) > 0 {
		buf.WriteString(cmd.Long + "\n\n")
	}

	if err := printOptions(buf, cmd); err != nil {
		return err
	}

	if len(cmd.Example) > 0 {
		buf.WriteString("### Examples\n\n")
		buf.WriteString(fmt.Sprintf("```sh\n%s\n```\n\n", cmd.Example))
	}

	if cmd.HasSubCommands() {
		buf.WriteString("### Commands\n\n")

		for _, child := range cmd.Commands() {
			if child == nil || !child.IsAvailableCommand() || child.IsAdditionalHelpTopicCommand() {
				continue
			}
			cname := name + " " + child.Name()
			link := cname + markdownExtension
			link = strings.ReplaceAll(link, " ", "-")
			buf.WriteString(fmt.Sprintf("* [%s](%s)\t - %s\n", cname, link, child.Short))
		}
		buf.WriteString("\n")
	}

	_, err := buf.WriteTo(w)
	return err
}

// CustomGenMarkdownTree generates markdown documentation for the entire command tree.
func CustomGenMarkdownTree(cmd *cobra.Command, dir string) error {
	for _, c := range cmd.Commands() {
		if !c.IsAvailableCommand() || c.IsAdditionalHelpTopicCommand() {
			continue
		}
		if err := CustomGenMarkdownTree(c, dir); err != nil {
			return err
		}
	}
	basename := strings.ReplaceAll(cmd.CommandPath(), " ", "-") + markdownExtension
	filename := filepath.Join(dir, basename)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)

	if err := CustomGenMarkdown(cmd, f); err != nil {
		return err
	}
	return nil
}

func main() {
	err := CustomGenMarkdownTree(cmd.RootCmd, "./doc")
	if err != nil {
		log.Fatal(err)
	}
}
