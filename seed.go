package main

import (
	"context"
	"fmt"
	"math/rand"
	"github.com/jackc/pgx/v5"
	"github.com/uber/h3-go/v4"
)

func main() {
	connStr := "postgresql://postgres.wkvbeyhnizxsftptbgtw:jayamalaKM@1982@aws-1-ap-southeast-2.pooler.supabase.com:6543/postgres"
	conn, _ := pgx.Connect(context.Background(), connStr)
	defer conn.Close(context.Background())

	// 1. Clear old Warsaw data first
	conn.Exec(context.Background(), "DELETE FROM drivers")

	fmt.Println("Seeding 1,000 drivers in Bengaluru...")

	for i := 1; i <= 1000; i++ {
		// Bengaluru Ranges
		lng := 77.50 + rand.Float64()*(77.75-77.50) // X
		lat := 12.85 + rand.Float64()*(13.05-12.85) // Y

		h3Cell, _ := h3.LatLngToCell(h3.NewLatLng(lat, lng), 8)
		driverID := fmt.Sprintf("blr_driver_%d", i)

		// REMEMBER: ST_MakePoint(Longitude, Latitude)
		_, err := conn.Exec(context.Background(), 
			"INSERT INTO drivers (id, location, h3_index) VALUES ($1, ST_SetSRID(ST_MakePoint($2, $3), 4326)::geography, $4)",
			driverID, lng, lat, h3Cell.String())
		
		if err != nil {
			fmt.Println("Error inserting:", err)
		}
	}
	fmt.Println("Successfully seeded Bengaluru!")
}