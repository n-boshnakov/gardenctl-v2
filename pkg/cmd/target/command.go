/*
SPDX-FileCopyrightText: 2021 SAP SE or an SAP affiliate company and Gardener contributors

SPDX-License-Identifier: Apache-2.0
*/

package target

import (
	"errors"
	"fmt"

	"github.com/gardener/gardenctl-v2/internal/util"
	commonTarget "github.com/gardener/gardenctl-v2/pkg/cmd/common/target"
	"github.com/gardener/gardenctl-v2/pkg/cmd/target/drop"
	"github.com/gardener/gardenctl-v2/pkg/target"

	"github.com/spf13/cobra"
)

// NewCommand returns a new target command.
func NewCommand(f util.Factory, o *Options, targetProvider *target.DynamicTargetProvider) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "target",
		Short: "Set scope for next operations, e.g. \"gardenctl target garden garden_name\" to target garden with name of garden_name",
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			suggestions, err := validArgsFunction(f, o, args, toComplete)
			if err != nil {
				fmt.Fprintln(o.IOStreams.ErrOut, err.Error())
				return nil, cobra.ShellCompDirectiveNoFileComp
			}

			return util.FilterStringsByPrefix(toComplete, suggestions), cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := o.Complete(f, cmd, args, targetProvider); err != nil {
				return fmt.Errorf("failed to complete command options: %w", err)
			}

			if err := o.Validate(); err != nil {
				return err
			}

			return runCommand(f, o)
		},
	}

	ioStreams := util.NewIOStreams()

	cmd.AddCommand(drop.NewCommand(f, drop.NewOptions(ioStreams), targetProvider))

	return cmd
}

func runCommand(f util.Factory, o *Options) error {
	manager, err := f.Manager()
	if err != nil {
		return err
	}

	ctx := f.Context()

	switch o.Kind {
	case commonTarget.TargetKindGarden:
		err = manager.TargetGarden(o.TargetName)
	case commonTarget.TargetKindProject:
		err = manager.TargetProject(ctx, o.TargetName)
	case commonTarget.TargetKindSeed:
		err = manager.TargetSeed(ctx, o.TargetName)
	case commonTarget.TargetKindShoot:
		err = manager.TargetShoot(ctx, o.TargetName)
	default:
		err = errors.New("invalid kind")
	}

	if err != nil {
		return err
	}

	fmt.Fprintf(o.IOStreams.Out, "Successfully targeted %s %q\n", o.Kind, o.TargetName)

	return nil
}
