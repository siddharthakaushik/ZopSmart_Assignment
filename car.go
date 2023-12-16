package main

import "gorm.io/gorm"

type Car struct {
   gorm.DB
   Brand  string
   Model  string
   Status string // "In Garage", "In Repair", "Completed", "Left Garage"
}
