package main

import (
	"fmt"

	"github.com/Thursdaylam/bom/pkg/node"
)

func main() {
	fmt.Println("Simulation starting")
	barcode()
	manual()
	rfid()
	fmt.Println("Simulation over")
}
func barcode() {

	var barcodeProb = make(map[string]float64)

	barcodeProb["findInPlant"] = 0.863
	barcodeProb["foundInPlantSearch"] = 0.972
	barcodeProb["findInSite"] = 0.847
	barcodeProb["foundInSiteSearch"] = 0.968
	barcodeProb["correctPiece"] = 0.971
	barcodeProb["foundWrongPiece"] = 0.97

	var barcodeCost = make(map[string]float64)

	barcodeCost["findInPlant"] = -0.77727
	barcodeCost["searchInPlant"] = -4.43127
	barcodeCost["replacementCost"] = -426

	barcodeCost["findInSite"] = -0.56127
	barcodeCost["foundInSite"] = 0
	barcodeCost["notFoundInSite"] = -7.13509

	barcodeCost["siteReplaceTransport"] = -571.455
	barcodeCost["siteTransport"] = -145.455
	barcodeCost["correctPiece"] = 0
	barcodeCost["wrongPiece"] = -2.70382

	barcodeCost["foundTransport"] = -145.455
	barcodeCost["missingReplaceTransport"] = -571.455

	average := createPCTree(barcodeProb, barcodeCost)

	fmt.Println("Average for barcode =", average)
}

func rfid() {

	var rfidProb = make(map[string]float64)

	rfidProb["findInPlant"] = 0.956
	rfidProb["foundInPlantSearch"] = 0.987
	rfidProb["findInSite"] = 0.962
	rfidProb["foundInSiteSearch"] = 0.978
	rfidProb["correctPiece"] = 0.991
	rfidProb["foundWrongPiece"] = 0.983

	var rfidCost = make(map[string]float64)

	rfidCost["findInPlant"] = -0.4211
	rfidCost["searchInPlant"] = -1.5813
	rfidCost["replacementCost"] = -426

	rfidCost["findInSite"] = -0.4271
	rfidCost["foundInSite"] = 0
	rfidCost["notFoundInSite"] = -2.3455

	rfidCost["siteReplaceTransport"] = -571.455
	rfidCost["siteTransport"] = -145.455
	rfidCost["correctPiece"] = 0
	rfidCost["wrongPiece"] = -0.7642

	rfidCost["foundTransport"] = -145.455
	rfidCost["missingReplaceTransport"] = -571.455

	average := createPCTree(rfidProb, rfidCost)

	fmt.Println("Average for rfid =", average)
}

func manual() {

	var manualProb = make(map[string]float64)

	manualProb["findInPlant"] = 0.712
	manualProb["foundInPlantSearch"] = 0.968
	manualProb["findInSite"] = 0.781
	manualProb["foundInSiteSearch"] = 0.965
	manualProb["correctPiece"] = 0.971
	manualProb["foundWrongPiece"] = 0.93

	var manualCost = make(map[string]float64)

	manualCost["findInPlant"] = -6.7871
	manualCost["searchInPlant"] = -9.598
	manualCost["replacementCost"] = -426

	manualCost["findInSite"] = -6.752
	manualCost["foundInSite"] = 0
	manualCost["notFoundInSite"] = -15.93

	manualCost["siteReplaceTransport"] = -571.455
	manualCost["siteTransport"] = -145.455
	manualCost["correctPiece"] = 0
	manualCost["wrongPiece"] = -6.334

	manualCost["foundTransport"] = -145.455
	manualCost["missingReplaceTransport"] = -571.455

	average := createPCTree(manualProb, manualCost)

	fmt.Println("Average for manual =", average)
}

func createPCTree(data map[string]float64, cost map[string]float64) float64 {
	findInPlant := node.Node{Value: cost["findInPlant"], Left: nil, Right: nil, Prob: data["findInPlant"]}

	foundInPlant := node.Node{Value: cost["findInSite"], Left: nil, Right: nil, Prob: data["findInSite"]}
	notFoundInPlant := node.Node{Value: cost["searchInPlant"], Left: nil, Right: nil, Prob: data["foundInPlantSearch"]}

	findInPlant.Left = &foundInPlant
	findInPlant.Right = &notFoundInPlant

	plantReplacement := node.Node{Value: cost["replacementCost"], Left: &foundInPlant, Right: &foundInPlant, Prob: 1.0}
	foundInSite := node.Node{Value: cost["foundInSite"], Left: nil, Right: nil, Prob: data["correctPiece"]}
	notFoundInSite := node.Node{Value: cost["notFoundInSite"], Left: nil, Right: nil, Prob: data["foundInSiteSearch"]}

	plantReplacement.Left = &findInPlant
	plantReplacement.Right = &findInPlant
	foundInPlant.Left = &foundInSite
	foundInPlant.Right = &notFoundInSite
	notFoundInPlant.Left = &foundInPlant
	notFoundInPlant.Right = &plantReplacement

	plantReplaceTransport := node.Node{Value: cost["siteReplaceTransport"], Left: nil, Right: nil, Prob: 1}
	foundInSiteExtendedSearch := node.Node{Value: cost["siteTransport"], Left: nil, Right: nil, Prob: 1}
	complete := node.Node{Value: cost["correctPiece"], Left: nil, Right: nil, Prob: 1}
	wrongPieceInSite := node.Node{Value: cost["wrongPiece"], Left: nil, Right: nil, Prob: data["foundWrongPiece"]}

	foundInSite.Left = &complete
	foundInSite.Right = &wrongPieceInSite
	notFoundInSite.Left = &foundInSiteExtendedSearch
	notFoundInSite.Right = &plantReplaceTransport

	plantReplaceTransport.Left = &findInPlant
	plantReplaceTransport.Right = &findInPlant
	foundInSiteExtendedSearch.Left = &foundInSite
	foundInSiteExtendedSearch.Right = &foundInSite

	foundMissing := node.Node{Value: cost["foundTransport"], Left: nil, Right: nil, Prob: 1}
	replaceTransportMissing := node.Node{Value: cost["missingReplaceTransport"], Left: nil, Right: nil, Prob: 1}

	wrongPieceInSite.Left = &foundMissing
	wrongPieceInSite.Right = &replaceTransportMissing

	foundMissing.Left = &complete
	foundMissing.Right = &complete
	replaceTransportMissing.Left = &findInPlant
	replaceTransportMissing.Right = &findInPlant

	value := node.EvaluateWorker(findInPlant, 1000000)

	return value

}
