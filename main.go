package main
import (
        "net/http"

        "github.com/gin-gonic/gin"
)
type counter struct{
        ID string `json:"id"`
        Name string `json:"name"`
        Value int `json:"value"`
}

var counters = []counter{
        {ID:"0", Name: "Zero", Value:0},
        {ID:"1", Name: "One", Value:0},
}
func main(){
        gin.SetMode(gin.DebugMode)
        router := gin.Default()
        router.GET("/counters", getCounters)
        router.PUT("/counters/:id", updateCounter)

        router.GET("/", func(c *gin.Context) {
                c.File("./public/index.html")
        })

        router.Run(":8080")

}


func getCounters(c *gin.Context){
        c.IndentedJSON(http.StatusOK, counters)
}

func updateCounter(c *gin.Context){
        id := c.Param("id")


        for i, count := range counters{
                if count.ID == id{
                        counters[i].Value  ++

                        c.JSON(http.StatusOK, gin.H{"message": "Counter updated"})
                        return
                }
        }

        c.JSON(http.StatusNotFound, gin.H{"message": "Counter not found"})
}
