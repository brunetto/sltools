package sla

var (
	// Regexp string for all_the_fishes
	regStringAllFishes = `^(\d{3})\s+` +			// GROUP 1: Z eg 001
								`(\d{3})\s+` + 		// GROUP 2: n eg 001
								`(\S+)\s+`+ 		// GROUP 3: binary_ids Z001n001idsa12550b2550
								`(\d+)\s+` + 		// GROUP 4: sys_time eg 0
								`(\d+\.\d+)\s+`+ 	// GROUP 5: phys_time [Myr] 0.0
								`(\S+)\s+`+ 		// GROUP 6: objects_ids eg  2550|12550
								`(\S)\s+`+ 			// GROUP 7: hardflag eg H 
								`(\S+)\s+`+ 		// GROUP 8: types eg ns++|ns++
								`(\S+\.\S+)\s+`+ 	// GROUP 9: masse[0] eg 10.3837569427
								`(\S+\.\S+)\s+`+	// GROUP 10: mass[1] eg 9.2141789593
								`(\S+\.\S+)\s+`+	// GROUP 11: sma eg  3.6333e-05
								`(\S+\.\S+)\s+`+	// GROUP 12: sma eg  4.6156152408e-06
								`(\S+\.*\S*)`		// GROUP 13: sma eg   0.680846
													// NOTE: Maybe ecc is zero...

	regStringDBHAll = `^(\d{3})\,\s+` +				// GROUP 1: Z eg 001
								`(\d{1,3})\,\s+` + 	// GROUP 2: n eg 001
								`(\S+)\,\s+`+ 		// GROUP 3: binary_ids Z001n001idsa12550b2550
								`(\d+)\,\s+` + 		// GROUP 4: sys_time eg 0
								`(\d+\.\d+)\,\s+`+ 	// GROUP 5: phys_time [Myr] 0.0
								`(\S+)\,\s+`+ 		// GROUP 6: objects_ids eg  2550|12550
								`(\S)\,\s+`+ 		// GROUP 7: hardflag eg H 
								`(\S+)\,\s+`+ 		// GROUP 8: types eg ns++|ns++
								`(\S+\.\S+)\,\s+`+ 	// GROUP 9: masse[0] eg 10.3837569427
								`(\S+\.\S+)\,\s+`+	// GROUP 10: mass[1] eg 9.2141789593
								`(\S+\.\S+)\,\s+`+	// GROUP 11: sma eg  3.6333e-05
								`(\S+\.\S+)\,\s+`+	// GROUP 12: sma eg  4.6156152408e-06
								`(\S+\.*\S*)`		// GROUP 13: sma eg   0.680846
													// NOTE: Maybe ecc is zero...
	inPath string
	inFile string
)