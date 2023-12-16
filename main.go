package main

import (
   "fmt"
   "net/http"
   "sync"
   "time"

   "github.com/gin-gonic/gin"
   "github.com/gorilla/websocket"
   "gorm.io/driver/sqlite"
   "gorm.io/gorm"
)

var (
   db            *gorm.DB
   upgrader      = websocket.Upgrader{}
   clientSockets = make(map[*websocket.Conn]bool)
   clientMutex   = sync.Mutex{}
)

func initDB() {
   var err error
   db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
   if err != nil {
      panic("failed to connect database")
   }

   db.AutoMigrate(&Car{})
}

func main() {
   r := gin.Default()

   initDB()

   // Routes
   r.POST("/cars", createCar)
   r.GET("/cars", getCarList)
   r.PUT("/cars/:id", updateCar)
   r.DELETE("/cars/:id", deleteCar)

   // WebSocket endpoint
   r.GET("/ws", handleWebSocket)

   // Run the server
   go handleWebSocketBroadcast()
   r.Run(":8080")
}

func createCar(c *gin.Context) {
   var car Car
   if err := c.ShouldBindJSON(&car); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
   }

   db.Create(&car)
   c.JSON(http.StatusCreated, car)
}

func getCarList(c *gin.Context) {
   var cars []Car
   db.Find(&cars)
   c.JSON(http.StatusOK, cars)
}

func updateCar(c *gin.Context) {
   id := c.Params.ByName("id")
   var car Car

   if err := db.Where("id = ?", id).First(&car).Error; err != nil {
      c.AbortWithStatus(http.StatusNotFound)
      return
   }

   var updatedCar Car
   if err := c.ShouldBindJSON(&updatedCar); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
      return
   }

   car.Status = updatedCar.Status
   db.Save(&car)

   c.JSON(http.StatusOK, car)
}

func deleteCar(c *gin.Context) {
   id := c.Params.ByName("id")
   var car Car

   if err := db.Where("id = ?", id).First(&car).Error; err != nil {
      c.AbortWithStatus(http.StatusNotFound)
      return
   }

   db.Delete(&car)
   c.JSON(http.StatusOK, gin.H{"id " + id: "deleted"})
}

func handleWebSocket(c *gin.Context) {
   conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
   if err != nil {
      fmt.Println(err)
      return
   }

   clientMutex.Lock()
   clientSockets[conn] = true
   clientMutex.Unlock()

   defer func() {
      clientMutex.Lock()
      delete(clientSockets, conn)
      clientMutex.Unlock()
      conn.Close()
   }()

   for {
      messageType, p, err := conn.ReadMessage()
      if err != nil {
         return
      }

      clientMutex.Lock()
      for client := range clientSockets {
         err := client.WriteMessage(messageType, p)
         if err != nil {
            fmt.Println(err)
            client.Close()
            delete(clientSockets, client)
         }
      }
      clientMutex.Unlock()
   }
}

func handleWebSocketBroadcast() {
   for {
      var cars []Car
      db.Find(&cars)

      clientMutex.Lock()
      for client := range clientSockets {
         err := client.WriteJSON(cars)
         if err != nil {
            fmt.Println(err)
            client.Close()
            delete(clientSockets, client)
         }
      }
      clientMutex.Unlock()

      // Broadcast every 5 seconds (adjust as needed)
      <-time.After(5 * time.Second)
   }
}
