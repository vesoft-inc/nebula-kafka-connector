/* Copyright (c) 2024 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package printer

import (
	"fmt"
	"log"

	nebula "github.com/vesoft-inc/nebula-ng-tools/golang"
)

type ExtraInfoPrinter struct{}

func (ExtraInfoPrinter) PrintMutationInfo(extraInfo nebula.ExtraInfo) {
	affectedNodes, ok := extraInfo.AffectedNodes()
	if ok != nil {
		return
	}
	affectedForwardEdges, okFor := extraInfo.AffectedForwardEdges()
	affectedReverseEdges, okRev := extraInfo.AffectedReverseEdges()
	if okFor != nil || okRev != nil {
		log.Fatal("affected_nodes, affected_forward_edges and affected_reverse_edges should be sent together.")
	}
	fmt.Println()
	fmt.Printf("Affected: %d nodes, %d forward edges, %d reverse edges\n",
		affectedNodes, affectedForwardEdges, affectedReverseEdges)
}
