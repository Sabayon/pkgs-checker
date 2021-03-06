/*

Copyright (C) 2017-2019  Daniele Rondina <geaaru@sabayonlinux.org>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

*/
package pkglist

import (
	"github.com/spf13/cobra"
)

func NewPkglistCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "pkglist [command] [OPTIONS]",
		Short: "Manage pkglist files.",
		Args:  cobra.OnlyValidArgs,
	}

	cmd.AddCommand(
		newPkglistCreateCommand(),
		newPkglistIntersectCommand(),
		newPkglistShowCommand(),
	)

	return cmd
}
