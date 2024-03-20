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
	affectedNodes, ok := extraInfo.GetValueByName("affected_nodes")
	if ok != nil {
		return
	}
	affectedForwardEdges, okFor := extraInfo.GetValueByName("affected_forward_edges")
	affectedReverseEdges, okRev := extraInfo.GetValueByName("affected_reverse_edges")
	if okFor != nil || okRev != nil {
		log.Fatal("affected_nodes, affected_forward_edges and affected_reverse_edges should be sent together.")
	}
	fmt.Printf("Affected: %s nodes, %s forward edges, %s reverse edges\n",
		affectedNodes.String(), affectedForwardEdges.String(), affectedReverseEdges.String())
}
