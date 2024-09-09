// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General
// Public License for more details.
//
// You should have received a copy of the GNU General Public License along
// with this program.  If not, see <http://www.gnu.org/licenses/>.

// Transform regular expressions into finite state machines.
package main

import (
	"flag"
	"fmt"
	"github.com/oulinbao/regexinter/intersection"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Println(`Usage: regexinter regexp1 regexp2

EXAMPLE: regexinter "a+b" "a*b"`)
	}
	flag.Parse()

	if len(flag.Args()) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	result, _ := intersection.HasIntersection(flag.Arg(0), flag.Arg(1))
	fmt.Println(result)
}
