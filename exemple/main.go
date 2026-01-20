package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/LIRYC-IHU/hl7v3-aecg/hl7aecg"
	"github.com/LIRYC-IHU/hl7v3-aecg/hl7aecg/types"
)

func main() {
	fmt.Println("=== HL7 aECG Example - Replication of Reference File ===")
	fmt.Println()

	// 1. Create a new aECG instance
	fmt.Println("1. Creating aECG instance...")
	h := hl7aecg.NewHl7xml("./output")

	// 2. Initialize with ECG procedure code
	fmt.Println("2. Initializing document...")
	h.Initialize(types.CPT_CODE_ECG_Routine, types.CPT_OID, "CPT-4", "")

	// Set the global root ID (singleton pattern)
	h.HL7AEcg.SetRootID("755.3045256.2025923.103550", "annotatedEcg")

	// Set confidentiality code (Normal)
	h.AddConfidentialityCode(types.CONFIDENTIALITY_INVESTIGATOR_BLINDED)

	// Set reason code (Per Protocol)
	h.AddReasonCode(types.REASON_PER_PROTOCOL)

	// 3. Set document effective time
	fmt.Println("3. Setting effective time...")
	startTime := "20250923103550"
	endTime := "20250923103600"
	h.SetEffectiveTime(startTime, endTime, nil, nil)

	// 4. Configure subject demographics (matching reference file)
	fmt.Println("4. Configuring subject...")
	// Patient ID from reference: 25060897140
	// Gender: UN (Undifferentiated)
	// Birth time: empty
	// Race: 2131-1 (Other Race)
	h.SetSubjectDemographics(
		"",                            // name - empty in reference
		"25060897140",                 // patientID
		types.GENDER_UNDIFFERENTIATED, // gender - UN
		"",                            // birthTime - empty in reference
		types.RaceCode("2131-1"),      // race code - Other Race
	)

	// Set the subject ID with extension
	subject := &h.HL7AEcg.ComponentOf.TimepointEvent.ComponentOf.SubjectAssignment.Subject.TrialSubject
	subject.ID.SetID("", "trialSubject")

	// 5. Configure clinical trial
	fmt.Println("5. Configuring clinical trial...")
	clinicalTrial := &h.HL7AEcg.ComponentOf.TimepointEvent.ComponentOf.SubjectAssignment.ComponentOf.ClinicalTrial
	clinicalTrial.SetID("", "")

	// Set Activity time
	clinicalTrial.SetActivityTime("432432", "432432")

	// 6. Configure location
	fmt.Println("6. Configuring location...")
	h.SetLocation(
		"trialSite",
		"",
		"ID302 - 3e EST - 1463",
		"",
		"",
		"",
	)

	// 7. Configure responsible party
	fmt.Println("7. Configuring responsible party...")
	h.SetResponsibleParty(
		"", // Use global root ID
		"trialInvestigator",
		"", // No prefix
		"", // No given name
		"", // No family name
		"", // No suffix - will create empty <name/> element
	)

	// 8. Generate simulated ECG data for 12 leads
	fmt.Println("8. Generating ECG data...")
	sampleRate := 100.0 // 500 Hz
	duration := 10.0    // 10 seconds
	numSamples := int(sampleRate * duration)

	// Create simulated data for all 12 leads
	leads := generateSimulatedECG(numSamples)

	// 9. Add rhythm series with series ID extension
	fmt.Println("9. Adding rhythm series...")
	startTimeMs := startTime
	endTimeMs := endTime

	h.AddRhythmSeries(
		startTimeMs,
		endTimeMs,
		types.BoolPtr(true),
		types.BoolPtr(false),
		sampleRate,
		leads,
		0.0, // origin
		5.0, // scale (5 µV per digit)
	)
	// Update series code with additional attributes

	if err := h.SetSeriesCode(
		types.RHYTHM_CODE,
		types.HL7_ActCode_OID,
		"ActCode",
		"Rhythm Waveforms",
	); err != nil {
		log.Printf("Warning: Could not set series code: %v\n", err)
	}

	// 10. Add device information (matching reference file)
	fmt.Println("10. Adding device information...")
	// Reference device: BeneHeart R12
	deviceModelName := "BeneHeart R12"
	softwareName := "02.10.00"
	manufacturerName := "(C) Shenzhen Mindray Bio-Medical Electronics Co., Ltd. All rights reserved."
	serialNumber := "FN-3A045256"

	if len(h.HL7AEcg.Component) > 0 {
		lastSeriesIdx := len(h.HL7AEcg.Component) - 1
		h.HL7AEcg.Component[lastSeriesIdx].Series.Author = &types.Author{
			SeriesAuthor: types.SeriesAuthor{
				ManufacturedSeriesDevice: types.ManufacturedSeriesDevice{
					ManufacturerModelName: &deviceModelName,
					SoftwareName:          &softwareName,
					SerialNumber:          &serialNumber,
				},
				ManufacturerOrganization: &types.ManufacturerOrganization{
					Name: &manufacturerName,
				},
			},
		}
	}

	// 11. Add secondary performer (technician)
	fmt.Println("11. Adding secondary performer...")
	h.AddSecondaryPerformer(types.PERFORMER_ECG_TECHNICIAN, "", "", "")

	// 12. Add control variables (filters)
	fmt.Println("12. Adding control variables (filters)...")
	// Match the reference file filters:
	// - Low Pass Filter: 35 Hz
	// - High Pass Filter: 0.56 Hz
	// - Notch Filter: 50 Hz
	h.AddLowPassFilter("35", "Hz")
	h.AddHighPassFilter("0.56", "Hz")
	h.AddNotchFilter("50", "Hz")

	// 13. Generate representative beat data
	fmt.Println("13. Generating representative beat data...")
	repBeatLeads := generateRepresentativeBeat(500) // 500 samples = ~1 second at 500 Hz

	// 14. Add derived representative beat series
	fmt.Println("14. Adding derived representative beat series...")
	h.AddDerivedSeries(
		types.REPRESENTATIVE_BEAT_CODE,
		startTime,            // Same start time as rhythm in absolute time
		endTime,              // End time
		types.BoolPtr(true),  // Start inclusive
		types.BoolPtr(false), // End exclusive
		sampleRate,           // Same sample rate as rhythm
		repBeatLeads,         // Representative beat data
		0.0,                  // origin
		5.0,                  // scale
	)

	// Set derived series properties
	if err := h.SetDerivedSeriesCode(
		types.REPRESENTATIVE_BEAT_CODE,
		types.HL7_ActCode_OID,
		"ActCode",
		"Representative Beat Waveforms",
	); err != nil {
		log.Printf("Warning: Could not set derived series code: %v\n", err)
	}

	// Set derived series ID
	if err := h.SetDerivedSeriesID("", "derivedSeries"); err != nil {
		log.Printf("Warning: Could not set derived series ID: %v\n", err)
	}

	// 15. Display summary
	fmt.Println()
	fmt.Println("=== Document Summary ===")
	fmt.Printf("Document ID: %s\n", h.HL7AEcg.ID.Root)
	fmt.Printf("Extension: %s\n", h.HL7AEcg.ID.Extension)
	fmt.Printf("Code: %s (%s)\n", h.HL7AEcg.Code.Code, h.HL7AEcg.Code.DisplayName)
	fmt.Printf("Effective Time: %s - %s\n", h.HL7AEcg.EffectiveTime.Low.Value, h.HL7AEcg.EffectiveTime.High.Value)

	if h.HL7AEcg.ComponentOf != nil {
		subject := &h.HL7AEcg.ComponentOf.TimepointEvent.ComponentOf.SubjectAssignment.Subject.TrialSubject
		if subject.ID != nil {
			fmt.Printf("\nSubject ID: %s (extension: %s)\n", subject.ID.Root, subject.ID.Extension)
			if subject.SubjectDemographicPerson != nil {
				if subject.SubjectDemographicPerson.PatientID != "" {
					fmt.Printf("Patient ID: %s\n", subject.SubjectDemographicPerson.PatientID)
				}
				if subject.SubjectDemographicPerson.AdministrativeGenderCode != nil {
					fmt.Printf("Gender: %s\n", subject.SubjectDemographicPerson.AdministrativeGenderCode.Code)
				}
			}
		}
	}

	fmt.Printf("\nNumber of series: %d\n", len(h.HL7AEcg.Component))
	if len(h.HL7AEcg.Component) > 0 {
		series := &h.HL7AEcg.Component[0].Series
		fmt.Printf("Series ID: %s (extension: %s)\n", series.ID.Root, series.ID.Extension)
		fmt.Printf("Series Type: %s\n", series.Code.Code)
		fmt.Printf("Number of sequences: %d\n", len(series.Component[0].SequenceSet.Component))

		if series.Author != nil {
			fmt.Printf("Device Model: %s\n", *series.Author.SeriesAuthor.ManufacturedSeriesDevice.ManufacturerModelName)
			fmt.Printf("Software: %s\n", *series.Author.SeriesAuthor.ManufacturedSeriesDevice.SoftwareName)
		}

		fmt.Printf("Secondary Performers: %d\n", len(series.SecondaryPerformer))
		fmt.Printf("Control Variables: %d\n", len(series.ControlVariable))

		// Show derived series info
		if len(series.Derivation) > 0 {
			fmt.Printf("\nNumber of derived series: %d\n", len(series.Derivation))
			for i, deriv := range series.Derivation {
				derivedSeries := deriv.DerivedSeries
				fmt.Printf("  Derived Series %d:\n", i+1)
				fmt.Printf("    ID Extension: %s\n", derivedSeries.ID.Extension)
				fmt.Printf("    Code: %s\n", derivedSeries.Code.Code)

				if len(derivedSeries.Component) > 0 {
					seqSet := derivedSeries.Component[0].SequenceSet
					fmt.Printf("    Number of sequences: %d\n", len(seqSet.Component))

					// Check time sequence type
					if len(seqSet.Component) > 0 {
						timeSeq := seqSet.Component[0].Sequence
						if timeSeq.Code.Time != nil {
							fmt.Printf("    Time sequence type: %s\n", timeSeq.Code.Time.Code)
						}

						// Check if GLIST_PQ
						if timeSeq.Value != nil && timeSeq.Value.XsiType == "GLIST_PQ" {
							if glistPq, ok := timeSeq.Value.Typed.(*types.GLIST_PQ); ok {
								fmt.Printf("    Time starts at: %s %s\n", glistPq.Head.Value, glistPq.Head.Unit)
								fmt.Printf("    Time increment: %s %s\n", glistPq.Increment.Value, glistPq.Increment.Unit)
							}
						}
					}
				}
			}
		}
	}

	// =========================================================================
	// Story U2: Add Annotations to Series
	// =========================================================================
	fmt.Println()
	fmt.Println("=== Story U2: Adding Annotations ===")
	fmt.Println()

	// Get the rhythm series to add annotations
	rhythmSeries := &h.HL7AEcg.Component[0].Series

	// Initialize annotation set with activity time
	annSet := rhythmSeries.InitAnnotationSet("20250923103600")
	fmt.Println("✓ Initialized annotation set")

	// Add global annotations (without supportingROI)
	// 1. Heart Rate
	annSet.AddHeartRate(57)
	fmt.Println("  Added annotation 1: Heart Rate = 57 bpm")

	// 2. PR Interval
	annSet.AddPRInterval(192)
	fmt.Println("  Added annotation 2: PR Interval = 192 ms")

	// 3. QRS Duration
	annSet.AddQRSDuration(88)
	fmt.Println("  Added annotation 3: QRS Duration = 88 ms")

	// 4. QT Interval
	annSet.AddQTInterval(418)
	fmt.Println("  Added annotation 4: QT Interval = 418 ms")

	// 5. QTc Interval with nested correction formula
	qtcIdx := annSet.AddQTcInterval(0) // Parent has no value, only nested
	qtcAnn := annSet.GetAnnotation(qtcIdx)
	if qtcAnn != nil {
		qtcAnn.AddNestedAnnotationWithCodeSystemName("MINDRAY_ECG_TIME_PD_QTcH", "MINDRAY", 413, "ms")
	}
	fmt.Println("  Added annotation 5: QTc with nested MINDRAY correction = 413 ms")

	// 6-8. Axis measurements (vendor-specific MINDRAY codes)
	annSet.AddAnnotationWithCodeSystemName("MINDRAY_ECG_P_AXIS", "MINDRAY", 60, "deg")
	fmt.Println("  Added annotation 6: P Axis = 60 deg")

	annSet.AddAnnotationWithCodeSystemName("MINDRAY_ECG_QRS_AXIS", "MINDRAY", 64, "deg")
	fmt.Println("  Added annotation 7: QRS Axis = 64 deg")

	annSet.AddAnnotationWithCodeSystemName("MINDRAY_ECG_T_AXIS", "MINDRAY", 11, "deg")
	fmt.Println("  Added annotation 8: T Axis = 11 deg")

	// 9-11. Amplitude measurements
	annSet.AddAnnotationWithCodeSystemName("MINDRAY_ECG_RV5", "MINDRAY", 1.008, "mV")
	fmt.Println("  Added annotation 9: RV5 = 1.008 mV")

	annSet.AddAnnotationWithCodeSystemName("MINDRAY_ECG_SV1", "MINDRAY", 0.565, "mV")
	fmt.Println("  Added annotation 10: SV1 = 0.565 mV")

	annSet.AddAnnotationWithCodeSystemName("MINDRAY_ECG_RV5_PLUS_SV1", "MINDRAY", 1.573, "mV")
	fmt.Println("  Added annotation 11: RV5+SV1 = 1.573 mV")

	// 12. Lead-specific annotation for Lead I with nested measurements
	leadIdx := annSet.AddLeadAnnotation("MDC_ECG_LEAD_I", "MINDRAY_MEASUREMENT_MATRIX", "MINDRAY_MEASUREMENT_MATRIX", "MINDRAY")
	leadAnn := annSet.GetAnnotation(leadIdx)
	if leadAnn != nil {
		leadAnn.AddNestedAnnotationWithCodeSystemName("MINDRAY_P_ONSET", "MINDRAY", 234, "ms")
		leadAnn.AddNestedAnnotationWithCodeSystemName("MINDRAY_P_DUR", "MINDRAY", 110, "ms")
		leadAnn.AddNestedAnnotationWithCodeSystemName("MINDRAY_QRS_ONSET", "MINDRAY", 428, "ms")
		leadAnn.AddNestedAnnotationWithCodeSystemName("MINDRAY_QRS_DURATION", "MINDRAY", 74, "ms")
		leadAnn.AddNestedAnnotationWithCodeSystemName("MINDRAY_Q_DUR", "MINDRAY", 20, "ms")
		leadAnn.AddNestedAnnotationWithCodeSystemName("MINDRAY_R_DUR", "MINDRAY", 34, "ms")
		leadAnn.AddNestedAnnotationWithCodeSystemName("MINDRAY_S_DUR", "MINDRAY", 18, "ms")
		leadAnn.AddNestedAnnotationWithCodeSystemName("MINDRAY_R_AMP", "MINDRAY", 535, "uV")
		leadAnn.AddNestedAnnotationWithCodeSystemName("MINDRAY_S_AMP", "MINDRAY", -165, "uV")
	}
	fmt.Println("  Added annotation 12: Lead I measurements (9 nested annotations)")

	// 13. Annotation with supportingROI for V6 (without nested measurements)
	annSet.AddLeadAnnotation("MDC_ECG_LEAD_V6", "MINDRAY_MEASUREMENT_MATRIX", "MINDRAY_MEASUREMENT_MATRIX", "MINDRAY")
	fmt.Println("  Added annotation 13: Lead V6 measurement container with supportingROI")

	// 14. ECG Interpretation with nested text statements
	interpIdx := annSet.AddTextAnnotation("MDC_ECG_INTERPRETATION", "2.16.840.1.113883.6.24", "")
	interpAnn := annSet.GetAnnotation(interpIdx)
	if interpAnn != nil {
		interpAnn.AddNestedTextAnnotation("MDC_ECG_INTERPRETATION_STATEMENT", "2.16.840.1.113883.6.24", "Rythme sinusal avec ESA")
		interpAnn.AddNestedTextAnnotation("MDC_ECG_INTERPRETATION_STATEMENT", "2.16.840.1.113883.6.24", "--- Interprétation sans connaître le sexe/l'âge du patient ---")
		interpAnn.AddNestedTextAnnotation("MDC_ECG_INTERPRETATION_SUMMARY", "2.16.840.1.113883.6.24", "ECG limite")
		interpAnn.AddNestedTextAnnotation("MDC_ECG_INTERPRETATION_COMMENT", "2.16.840.1.113883.6.24", "Diagnostic non confirmé.")
	}
	fmt.Println("  Added annotation 14: ECG Interpretation with 4 nested text statements")

	fmt.Printf("\n✓ Total: 14 annotations added to series (24+ including nested)\n")

	// 15. Validate the document
	fmt.Println()
	fmt.Println("15. Validating document...")
	if err := h.Validate(); err != nil {
		log.Printf("WARNING: Validation errors detected:\n%v\n", err)
	} else {
		fmt.Println("✓ Document validated successfully!")
	}

	// 16. Export to XML
	fmt.Println()
	fmt.Println("16. Exporting XML...")
	if _, err := h.Test(); err != nil {
		log.Fatalf("Error exporting XML: %v", err)
	}

	fmt.Println()
	fmt.Println("=== Generation Complete! ===")
	fmt.Println("Output file: /tmp/hl7aecg_example.xml")
	fmt.Println("Reference file: 25060897140_23092025103550.xml")
	fmt.Println()
	fmt.Printf("Total ECG data: %d samples per lead\n", numSamples)
	fmt.Printf("Sample rate: %.0f Hz\n", sampleRate)
	fmt.Printf("Duration: %.1f seconds\n", duration)
	// =========================================================================
	// Story 1.2: Demonstrate Unmarshal API
	// =========================================================================
	fmt.Println()
	fmt.Println("=== Story 1.2: Unmarshal API Demo ===")
	fmt.Println()

	// Method 1: Unmarshal from byte slice
	fmt.Println("Method 1: h.Unmarshal([]byte) - Parse from byte slice")
	fmt.Println("-----------------------------------------------------")
	data, err := os.ReadFile("./25060897140_23092025103550.xml")
	if err != nil {
		fmt.Printf("Error reading reference file: %v\n", err)
		return
	}

	parsedDoc := hl7aecg.NewHl7xml("")
	if err := parsedDoc.Unmarshal(data); err != nil {
		fmt.Printf("Error unmarshalling: %v\n", err)
		return
	}

	fmt.Println("✓ Document successfully parsed from byte slice!")
	fmt.Printf("  Document ID: %s\n", parsedDoc.HL7AEcg.ID.Root)
	fmt.Printf("  Document Code: %s\n", parsedDoc.HL7AEcg.Code.Code)
	if parsedDoc.HL7AEcg.EffectiveTime != nil {
		fmt.Printf("  Effective Time: %s - %s\n",
			parsedDoc.HL7AEcg.EffectiveTime.Low.Value,
			parsedDoc.HL7AEcg.EffectiveTime.High.Value)
	}

	// Show series info
	if len(parsedDoc.HL7AEcg.Component) > 0 {
		series := &parsedDoc.HL7AEcg.Component[0].Series
		fmt.Printf("  Series Code: %s\n", series.Code.Code)
		if len(series.Component) > 0 {
			seqSet := series.Component[0].SequenceSet
			fmt.Printf("  Number of sequences: %d\n", len(seqSet.Component))

			// Show SequenceValue polymorphism in action
			if len(seqSet.Component) > 0 {
				timeSeq := seqSet.Component[0].Sequence
				if timeSeq.Value != nil {
					fmt.Printf("  Time sequence type: %s\n", timeSeq.Value.XsiType)
					if glistTs, ok := timeSeq.Value.Typed.(*types.GLIST_TS); ok {
						fmt.Printf("    Head: %s\n", glistTs.Head.Value)
						fmt.Printf("    Increment: %s %s\n", glistTs.Increment.Value, glistTs.Increment.Unit)
					}
				}
			}

			// Show a lead sequence with SLIST_PQ
			if len(seqSet.Component) > 1 {
				leadSeq := seqSet.Component[1].Sequence
				if leadSeq.Value != nil {
					fmt.Printf("  Lead sequence type: %s\n", leadSeq.Value.XsiType)
					if slistPq, ok := leadSeq.Value.Typed.(*types.SLIST_PQ); ok {
						fmt.Printf("    Origin: %s %s\n", slistPq.Origin.Value, slistPq.Origin.Unit)
						fmt.Printf("    Scale: %s %s\n", slistPq.Scale.Value, slistPq.Scale.Unit)
						// Show first few digits
						digits := slistPq.Digits
						if len(digits) > 50 {
							digits = digits[:50] + "..."
						}
						fmt.Printf("    Digits: %s\n", digits)
					}
				}
			}
		}
	}

	fmt.Println()

	// Method 2: UnmarshalFromReader
	fmt.Println("Method 2: h.UnmarshalFromReader(io.Reader) - Parse from reader")
	fmt.Println("--------------------------------------------------------------")

	file, err := os.Open("./25060897140_23092025103550.xml")
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	parsedDoc2 := hl7aecg.NewHl7xml("")
	if err := parsedDoc2.UnmarshalFromReader(file); err != nil {
		fmt.Printf("Error unmarshalling from reader: %v\n", err)
		return
	}

	fmt.Println("✓ Document successfully parsed from io.Reader!")
	fmt.Printf("  Document ID: %s\n", parsedDoc2.HL7AEcg.ID.Root)
	fmt.Printf("  Number of series: %d\n", len(parsedDoc2.HL7AEcg.Component))

	fmt.Println()

	// Method 3: Error handling demonstration
	fmt.Println("---------------------------------------------")
	fmt.Println("Method 3: Error handling with wrapped errors")
	fmt.Println("---------------------------------------------")
	invalidXML := []byte("<invalid>xml<broken")
	errDoc := hl7aecg.NewHl7xml("")
	if err := errDoc.Unmarshal(invalidXML); err != nil {
		fmt.Printf("✓ Error correctly returned: %v\n", err)
		// Show that errors.Unwrap works
		if unwrapped := errors.Unwrap(err); unwrapped != nil {
			fmt.Printf("  Unwrapped error: %v\n", unwrapped)
		}
	}

	// Method 4: Unmarshal and Validate
	fmt.Println("---------------------------------------------")
	fmt.Println("Method 4: h.UnmarshalAndValidate([]byte) - Parse and validate")
	fmt.Println("---------------------------------------------")
	parsedDoc3 := hl7aecg.NewHl7xml("")
	err = parsedDoc3.UnmarshalAndValidate(data)
	if err != nil {
		fmt.Printf("Error unmarshalling and validating: %v\n", err)
	}

	fmt.Println()
	fmt.Println("=== Story 1.2 Demo Complete ===")

	// =========================================================================
	// Story U1: Annotation Set Unmarshalling Demo
	// =========================================================================
	fmt.Println()
	fmt.Println("=== Story U1: Annotation Set Demo ===")
	fmt.Println()

	// Access annotations from the parsed document
	if len(parsedDoc.HL7AEcg.Component) > 0 {
		series := &parsedDoc.HL7AEcg.Component[0].Series

		// Check if SubjectOf (annotations) exists
		if len(series.SubjectOf) > 0 && series.SubjectOf[0].AnnotationSet != nil {
			annSet := series.SubjectOf[0].AnnotationSet

			// Display activity time
			if annSet.ActivityTime != nil {
				fmt.Printf("Annotation Activity Time: %s\n", annSet.ActivityTime.Value)
			}
			fmt.Println()

			// Display first 20 annotations
			fmt.Println("First 20 Annotations:")
			fmt.Println("---------------------")

			count := 0
			maxAnnotations := 20

			for _, comp := range annSet.Component {
				if count >= maxAnnotations {
					break
				}

				ann := &comp.Annotation

				// Check if this is a global annotation (no support/ROI)
				if ann.Support == nil {
					// Global annotation
					if ann.Code != nil {
						fmt.Printf("%2d. [GLOBAL] %s", count+1, ann.Code.Code)

						// Display value if present
						if ann.Value != nil {
							if val, ok := ann.Value.GetValueFloat(); ok {
								fmt.Printf(" = %.0f %s", val, ann.Value.GetValueUnit())
							} else if text, ok := ann.Value.GetText(); ok {
								fmt.Printf(" = %s", text)
							}
						}
						fmt.Println()

						// Check for nested annotations (like QTc with correction methods)
						for _, nestedComp := range ann.Component {
							nestedAnn := &nestedComp.Annotation
							if nestedAnn.Code != nil {
								count++
								if count >= maxAnnotations {
									break
								}
								fmt.Printf("%2d.   └─ [NESTED] %s", count+1, nestedAnn.Code.Code)
								if nestedAnn.Value != nil {
									if val, ok := nestedAnn.Value.GetValueFloat(); ok {
										fmt.Printf(" = %.0f %s", val, nestedAnn.Value.GetValueUnit())
									} else if text, ok := nestedAnn.Value.GetText(); ok {
										fmt.Printf(" = %s", text)
									}
								}
								fmt.Println()
							}
						}
					}
				} else {
					// Lead-specific annotation
					leadCode := ""
					if len(ann.Support.SupportingROI.Component) > 0 {
						leadCode = string(ann.Support.SupportingROI.Component[0].Boundary.Code.Code)
					}

					if ann.Code != nil {
						fmt.Printf("%2d. [LEAD: %s] %s\n", count+1, leadCode, ann.Code.Code)

						// Show first few nested measurements for lead-specific annotations
						nestedShown := 0
						for _, nestedComp := range ann.Component {
							if nestedShown >= 3 { // Show only first 3 nested annotations per lead
								fmt.Printf("       ... and %d more measurements\n", len(ann.Component)-nestedShown)
								break
							}
							nestedAnn := &nestedComp.Annotation
							if nestedAnn.Code != nil && nestedAnn.Value != nil {
								if val, ok := nestedAnn.Value.GetValueFloat(); ok {
									fmt.Printf("       └─ %s = %.0f %s\n",
										nestedAnn.Code.Code,
										val,
										nestedAnn.Value.GetValueUnit())
								} else if text, ok := nestedAnn.Value.GetText(); ok {
									fmt.Printf("       └─ %s = %s\n",
										nestedAnn.Code.Code,
										text)
								}
								nestedShown++
							}
						}
					}
				}
				count++
			}

			fmt.Println()
			fmt.Printf("Total annotations in file: %d\n", len(annSet.Component))

			// Demonstrate helper methods
			fmt.Println()
			fmt.Println("Using Helper Methods:")
			fmt.Println("--------------------")

			// Get heart rate
			hrAnn := annSet.GetAnnotationByCode("MDC_ECG_HEART_RATE")
			if hrAnn != nil {
				val, ok := hrAnn.GetValueFloat()
				if ok {
					fmt.Printf("Heart Rate: %.0f %s\n", val, hrAnn.GetValueUnit())
				}
			}

			// Get PR interval
			prAnn := annSet.GetAnnotationByCode("MDC_ECG_TIME_PD_PR")
			if prAnn != nil {
				val, ok := prAnn.GetValueFloat()
				if ok {
					fmt.Printf("PR Interval: %.0f %s\n", val, prAnn.GetValueUnit())
				}
			}

			// Get QRS duration
			qrsAnn := annSet.GetAnnotationByCode("MDC_ECG_TIME_PD_QRS")
			if qrsAnn != nil {
				val, ok := qrsAnn.GetValueFloat()
				if ok {
					fmt.Printf("QRS Duration: %.0f %s\n", val, qrsAnn.GetValueUnit())
				}
			}

			// Get QT interval
			qtAnn := annSet.GetAnnotationByCode("MDC_ECG_TIME_PD_QT")
			if qtAnn != nil {
				val, ok := qtAnn.GetValueFloat()
				if ok {
					fmt.Printf("QT Interval: %.0f %s\n", val, qtAnn.GetValueUnit())
				}
			}

			// Get QTc with nested annotation
			qtcAnn := annSet.GetAnnotationByCode("MDC_ECG_TIME_PD_QTc")
			if qtcAnn != nil {
				nestedQtc := qtcAnn.GetNestedAnnotationByCode("MINDRAY_ECG_TIME_PD_QTcH")
				if nestedQtc != nil {
					val, ok := nestedQtc.GetValueFloat()
					if ok {
						fmt.Printf("QTc (Hodges): %.0f %s\n", val, nestedQtc.GetValueUnit())
					}
				}
			}

			// Get Lead I annotations
			fmt.Println()
			fmt.Println("Lead I Measurements:")
			leadIAnn := annSet.GetLeadAnnotations("MDC_ECG_LEAD_I")
			if leadIAnn != nil {
				// Get R amplitude
				rAmp := leadIAnn.GetNestedAnnotationByCode("MINDRAY_R_AMP")
				if rAmp != nil {
					val, ok := rAmp.GetValueFloat()
					if ok {
						fmt.Printf("  R Amplitude: %.0f %s\n", val, rAmp.GetValueUnit())
					}
				}

				// Get S amplitude
				sAmp := leadIAnn.GetNestedAnnotationByCode("MINDRAY_S_AMP")
				if sAmp != nil {
					val, ok := sAmp.GetValueFloat()
					if ok {
						fmt.Printf("  S Amplitude: %.0f %s\n", val, sAmp.GetValueUnit())
					}
				}

				// Get P onset
				pOnset := leadIAnn.GetNestedAnnotationByCode("MINDRAY_P_ONSET")
				if pOnset != nil {
					val, ok := pOnset.GetValueFloat()
					if ok {
						fmt.Printf("  P Onset: %.0f %s\n", val, pOnset.GetValueUnit())
					}
				}
			}
		} else {
			fmt.Println("No annotations found in the parsed document.")
		}
	}

	fmt.Println()
	fmt.Println("=== Story U1 Demo Complete ===")
	// == End of main function ==
}

// generateSimulatedECG generates simulated ECG data for all 12 standard leads
// This creates realistic-looking ECG waveforms with QRS complexes
func generateSimulatedECG(numSamples int) map[types.LeadCode][]int {
	leads := make(map[types.LeadCode][]int)

	// Get standard 12 leads
	standardLeads := types.GetStandardLeads()

	// Simulation parameters (at 500 Hz)
	heartRate := 60.0                               // beats per minute
	samplesPerBeat := int(500.0 * 60.0 / heartRate) // ~500 samples per beat at 60 bpm

	// Generate simulated data for each lead with lead-specific amplitudes
	for idx, lead := range standardLeads {
		samples := make([]int, numSamples)

		// Lead-specific amplitude multipliers (approximate realistic values)
		amplitude := getLeadAmplitude(idx)
		baseline := getLeadBaseline(idx)

		for i := 0; i < numSamples; i++ {
			// Position within the cardiac cycle
			beatPos := i % samplesPerBeat

			// Generate ECG waveform components
			// P wave (atrial depolarization) - small positive deflection
			pWave := 0.0
			if beatPos >= 50 && beatPos < 100 {
				pWave = 15.0 * amplitude * float64(beatPos-50) / 50.0
			} else if beatPos >= 100 && beatPos < 150 {
				pWave = 15.0 * amplitude * float64(150-beatPos) / 50.0
			}

			// QRS complex (ventricular depolarization) - sharp spike
			qrs := 0.0
			if beatPos >= 200 && beatPos < 210 {
				// Q wave (small negative)
				qrs = -20.0 * amplitude
			} else if beatPos >= 210 && beatPos < 230 {
				// R wave (large positive)
				qrs = 400.0 * amplitude * float64(beatPos-210) / 20.0
			} else if beatPos >= 230 && beatPos < 250 {
				// S wave (negative)
				qrs = -100.0 * amplitude * float64(beatPos-230) / 20.0
			}

			// T wave (ventricular repolarization) - rounded positive deflection
			tWave := 0.0
			if beatPos >= 300 && beatPos < 350 {
				tWave = 50.0 * amplitude * float64(beatPos-300) / 50.0
			} else if beatPos >= 350 && beatPos < 400 {
				tWave = 50.0 * amplitude * float64(400-beatPos) / 50.0
			}

			// Add some random noise for realism
			noise := float64((i*7)%5 - 2) // deterministic "noise" for reproducibility

			// Combine all components
			value := baseline + pWave + qrs + tWave + noise
			samples[i] = int(value)
		}

		leads[lead] = samples
	}

	return leads
}

// getLeadAmplitude returns the amplitude multiplier for each lead
func getLeadAmplitude(leadIndex int) float64 {
	amplitudes := []float64{
		1.0, // Lead I
		1.5, // Lead II
		0.5, // Lead III
		0.8, // aVR
		0.6, // aVL
		1.2, // aVF
		0.4, // V1
		0.7, // V2
		1.3, // V3
		1.5, // V4
		1.4, // V5
		1.1, // V6
	}
	if leadIndex < len(amplitudes) {
		return amplitudes[leadIndex]
	}
	return 1.0
}

// getLeadBaseline returns the baseline voltage for each lead
func getLeadBaseline(leadIndex int) float64 {
	baselines := []float64{
		-40, // Lead I
		-30, // Lead II
		-20, // Lead III
		-35, // aVR
		-25, // aVL
		-30, // aVF
		-15, // V1
		-20, // V2
		-25, // V3
		-30, // V4
		-35, // V5
		-40, // V6
	}
	if leadIndex < len(baselines) {
		return baselines[leadIndex]
	}
	return -30.0
}

// generateRepresentativeBeat generates a single representative beat for all 12 leads.
//
// In a real implementation, this would:
//  1. Detect QRS complexes in rhythm data
//  2. Extract beats around each QRS
//  3. Align beats by fiducial point (e.g., R-peak)
//  4. Average aligned beats to get representative beat
//
// This simplified version creates a single idealized beat waveform.
func generateRepresentativeBeat(numSamples int) map[types.LeadCode][]int {
	leads := make(map[types.LeadCode][]int)
	standardLeads := types.GetStandardLeads()

	for idx, lead := range standardLeads {
		samples := make([]int, numSamples)
		amplitude := getLeadAmplitude(idx)
		baseline := getLeadBaseline(idx)

		for i := 0; i < numSamples; i++ {
			// Normalized position in beat (0.0 to 1.0)
			pos := float64(i) / float64(numSamples)

			// P wave (10-20% of cycle)
			pWave := 0.0
			if pos >= 0.1 && pos < 0.15 {
				pWave = 15.0 * amplitude * (pos - 0.1) / 0.05
			} else if pos >= 0.15 && pos < 0.2 {
				pWave = 15.0 * amplitude * (0.2 - pos) / 0.05
			}

			// QRS complex (40-45% of cycle)
			qrs := 0.0
			if pos >= 0.40 && pos < 0.42 {
				qrs = -20.0 * amplitude // Q wave
			} else if pos >= 0.42 && pos < 0.44 {
				qrs = 400.0 * amplitude * (pos - 0.42) / 0.02 // R wave up
			} else if pos >= 0.44 && pos < 0.46 {
				qrs = -100.0 * amplitude * (pos - 0.44) / 0.02 // S wave down
			}

			// T wave (60-80% of cycle)
			tWave := 0.0
			if pos >= 0.60 && pos < 0.70 {
				tWave = 50.0 * amplitude * (pos - 0.60) / 0.10
			} else if pos >= 0.70 && pos < 0.80 {
				tWave = 50.0 * amplitude * (0.80 - pos) / 0.10
			}

			// Combine components
			value := baseline + pWave + qrs + tWave
			samples[i] = int(value)
		}

		leads[lead] = samples
	}

	return leads
}
