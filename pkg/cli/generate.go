package cli

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"sigs.k8s.io/kubebuilder/v3/pkg/config"
	"sigs.k8s.io/kubebuilder/v3/pkg/plugin"
)

func (c CLI) newAlphaCommandGenerate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "Re-scaffold",
		Short: "Rescaffold a project",
		Long:  `upgrade projects by re-scaffolding`,
		Run:     func(cmd *cobra.Command, args []string) {},
	}
	// Register --project-version on the dynamically created command
	// so that it shows up in help and does not cause a parse error.
	cmd.Flags().String(projectVersionFlag, c.defaultProjectVersion.String(), "project version")

	// In case no plugin was resolved, instead of failing the construction of the CLI, fail the execution of
	// this subcommand. This allows the use of subcommands that do not require resolved plugins like help.
	if len(c.resolvedPlugins) == 0 {
		cmdErr(cmd, noResolvedPluginError{})
		return cmd
	}

	// Obtain the plugin keys and subcommands from the plugins that implement plugin.Init.
	subcommands := c.filterSubcommands(
		func(p plugin.Plugin) bool {
			_, isValid := p.(plugin.Generate)
			return isValid
		},
		func(p plugin.Plugin) plugin.Subcommand {
			return p.(plugin.Generate).GetGenerateSubcommand()
		},
	)

	// Verify that there is at least one remaining plugin.
	if len(subcommands) == 0 {
		cmdErr(cmd, noAvailablePluginError{"project initialization"})
		return cmd
	}

	c.applySubcommandHooks(cmd, subcommands, initErrorMsg, true)

	return cmd
}

// func (c CLI) getGenerateExamples() string {
// 	var sb strings.Builder
// 	for _, version := range c.getAvailableProjectVersions() {
// 		rendered := fmt.Sprintf(`  # Help for initializing a project with version %[2]s
//   %[1]s init --project-version=%[2]s -h

// `,
// 			c.commandName, version)
// 		sb.WriteString(rendered)
// 	}
// 	return strings.TrimSuffix(sb.String(), "\n\n")
// }

// func (c CLI) getVersions() (projectVersions []string) {
// 	versionSet := make(map[config.Version]struct{})
// 	for _, p := range c.plugins {
// 		// Only return versions of non-deprecated plugins.
// 		if _, isDeprecated := p.(plugin.Deprecated); !isDeprecated {
// 			for _, version := range p.SupportedProjectVersions() {
// 				versionSet[version] = struct{}{}
// 			}
// 		}
// 	}
// 	for version := range versionSet {
// 		projectVersions = append(projectVersions, strconv.Quote(version.String()))
// 	}
// 	sort.Strings(projectVersions)
// 	return projectVersions
// }