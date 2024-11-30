package main

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

const (
	BatteryCost      = 500
	PaintThinnerCost = 300
	DrainCleanerCost = 300
)

type Handler func(w http.ResponseWriter, r *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		fmt.Println("internal server error:", err)
		w.WriteHeader(503)
		w.Write([]byte(fmt.Sprintf("internal server error: %v", err)))
	}
}

func (c *Config) HomeHandler(w http.ResponseWriter, r *http.Request) error {
	data := make(map[string]interface{})

	cookie, err := r.Cookie("session_token")
	if err != nil {
		data["IsAuthorized"] = false
		return c.tpl.ExecuteTemplate(w, "index", data)
	}
	sessionID := cookie.Value
	sessionUUID, err := uuid.Parse(sessionID)
	if err != nil {
		data["IsAuthorized"] = false
		return c.tpl.ExecuteTemplate(w, "index", data)
	}

	if !c.SessionStore.ValidateSession(sessionUUID) {
		data["IsAuthorized"] = false
		return c.tpl.ExecuteTemplate(w, "index", data)
	}

	session, _ := c.SessionStore.GetSession(sessionUUID)

	// // create active connection to database
	// c.StorageToolStore.Connect(session.LoginID)

	storages, err := c.StorageToolStore.GetGroupStorage(session.LoginID)
	if err != nil {
		log.Fatalln(err)
	}

	propertyInventories := []map[string]interface{}{}
	for _, storage := range storages {
		storageData := map[string]interface{}{
			"ID":           storage.ID,
			"PropertyName": storage.StorageName,
			"SlabsClear":   storage.ClearSlabQty.Int32,
			"BlocksClear":  storage.ClearBlockQty.Int32,
			"SlabsCloudy":  storage.CloudySlabQty.Int32,
			"BlocksCloudy": storage.CloudyBlockQty.Int32,
		}
		propertyInventories = append(propertyInventories, storageData)
	}
	data["IsAuthorized"] = true
	data["PropertyInventories"] = propertyInventories

	return c.tpl.ExecuteTemplate(w, "index", data)
}

func (c *Config) StorageToolHandler(w http.ResponseWriter, r *http.Request) error {
	data := make(map[string]interface{})

	cookie, err := r.Cookie("session_token")
	if err != nil {
		data["IsAuthorized"] = false
	} else {
		sessionID := cookie.Value
		sessionUUID, err := uuid.Parse(sessionID)
		if err != nil {
			log.Println("Error parsing UUID", err)
			data["IsAuthorized"] = false
		} else if c.SessionStore.ValidateSession(sessionUUID) {
			session, _ := c.SessionStore.GetSession(sessionUUID)

			// c.StorageToolStore.Connect(session.LoginID)

			storages, err := c.StorageToolStore.GetGroupStorage(session.LoginID)
			if err != nil {
				log.Fatalln(err)
			}
			propertyInventories := []map[string]interface{}{}

			for _, storage := range storages {
				storageData := map[string]interface{}{
					"ID":           storage.ID,
					"PropertyName": storage.StorageName,
					"SlabsClear":   storage.ClearSlabQty.Int32,
					"BlocksClear":  storage.ClearBlockQty.Int32,
					"SlabsCloudy":  storage.CloudySlabQty.Int32,
					"BlocksCloudy": storage.CloudyBlockQty.Int32,
				}
				propertyInventories = append(propertyInventories, storageData)
			}
			data["IsAuthorized"] = true
			data["PropertyInventories"] = propertyInventories
		} else {
			data["IsAuthorized"] = false
		}
	}
	data["WsUrl"] = c.WsUrl
	buf := bytes.Buffer{}
	c.tpl.ExecuteTemplate(&buf, "storageTool", data)

	return c.tpl.ExecuteTemplate(w, "storageTool", data)
}

func (c *Config) CookToolHandler(w http.ResponseWriter, r *http.Request) error {
	recipes := []Recipe{
		{
			ID: 1,
			Ingredients: []Ingredient{
				{Name: "Lithium", Amount: 5},
				{Name: "Acetone", Amount: 4},
				{Name: "Sulfuric Acid", Amount: 6},
			},
			IsSelected: false,
		},
		{
			ID: 2,
			Ingredients: []Ingredient{
				{Name: "Lithium", Amount: 7},
				{Name: "Acetone", Amount: 4},
				{Name: "Sulfuric Acid", Amount: 6},
			},
			IsSelected: false,
		},
	}

	suppliesOnHand := []Ingredient{
		{Name: "Lithium", Amount: 0},
		{Name: "Acetone", Amount: 0},
		{Name: "Sulfuric Acid", Amount: 0},
	}

	data := map[string]interface{}{
		"Recipes":        recipes,
		"SuppliesOnHand": suppliesOnHand,
		"NumberOfCooks":  0,
		"Results": map[string]interface{}{
			"CarBatteryNeeded":   0,
			"PaintThinnerNeeded": 0,
			"DrainCleanerNeeded": 0,
			"PrecursorWeight":    0,
			"PrecursorCost":      0,
		},
	}

	return c.tpl.ExecuteTemplate(w, "cookTool", data)
}

func (c *Config) TruckToolHandler(w http.ResponseWriter, r *http.Request) error {
	return c.tpl.ExecuteTemplate(w, "truckTool", nil)
}

func (c *Config) CookCalculateHandler(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	selectedRecipes := make(map[int]map[string]int)
	for k, v := range r.Form {
		if strings.HasSuffix(k, "-checkbox") && v[0] == "on" {
			idStr := strings.TrimPrefix(strings.TrimSuffix(k, "-checkbox"), "recipe-")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				continue
			}

			selectedRecipes[id] = make(map[string]int)

			for key, val := range r.Form {
				if strings.HasPrefix(key, fmt.Sprintf("recipe-%d", id)) && !strings.HasSuffix(key, "checkbox") {
					ingredientName := strings.TrimPrefix(key, fmt.Sprintf("recipe-%d-", id))
					ingredientAmount, _ := strconv.Atoi(val[0])
					selectedRecipes[id][ingredientName] = ingredientAmount
				}
			}
		}
	}
	maxLithium, maxAcetone, maxSulfuric := 0, 0, 0
	for _, ingredientMap := range selectedRecipes {
		if ingredientMap["Lithium"] > maxLithium {
			maxLithium = ingredientMap["Lithium"]
		}

		if ingredientMap["Acetone"] > maxAcetone {
			maxAcetone = ingredientMap["Acetone"]
		}

		if ingredientMap["Sulfuric Acid"] > maxSulfuric {
			maxSulfuric = ingredientMap["Sulfuric Acid"]
		}
	}
	numCooksStr := r.FormValue("number-of-cooks")
	numCooks, _ := strconv.Atoi(numCooksStr)

	lithiumOnHand, _ := strconv.Atoi(r.FormValue("Lithium-on-hand"))
	acetoneOnHand, _ := strconv.Atoi(r.FormValue("Acetone-on-hand"))
	sulfuricOnHand, _ := strconv.Atoi(r.FormValue("Sulfuric Acid-on-hand"))

	requiredLithium := maxLithium*numCooks - lithiumOnHand
	requiredAcetone := maxAcetone*numCooks - acetoneOnHand
	requiredSulfuric := maxSulfuric*numCooks - sulfuricOnHand

	precursorWeight := requiredLithium + requiredAcetone + requiredSulfuric

	carBatteryNeeded := int(math.Ceil(float64(requiredLithium) / 10))
	paintThinnerNeeded := int(math.Ceil(float64(requiredAcetone) / 5))
	drainCleanerNeeded := int(math.Ceil(float64(requiredSulfuric) / 5))

	if carBatteryNeeded < 0 {
		carBatteryNeeded = 0
	}

	if paintThinnerNeeded < 0 {
		paintThinnerNeeded = 0
	}

	if drainCleanerNeeded < 0 {
		drainCleanerNeeded = 0
	}

	carBatteryCost := carBatteryNeeded * BatteryCost
	drainCleanerCost := drainCleanerNeeded * DrainCleanerCost
	paintThinnerCost := paintThinnerNeeded * PaintThinnerCost

	precursorCost := carBatteryCost + drainCleanerCost + paintThinnerCost

	data := map[string]interface{}{
		"CarBatteryNeeded":   carBatteryNeeded,
		"PaintThinnerNeeded": paintThinnerNeeded,
		"DrainCleanerNeeded": drainCleanerNeeded,
		"PrecursorWeight":    precursorWeight,
		"PrecursorCost":      precursorCost,
	}

	w.Header().Set("Content-Type", "text/html")

	err = c.tpl.ExecuteTemplate(w, "cookToolResults", data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) CookDeleteRecipeHandler(w http.ResponseWriter, r *http.Request) error {
	recipeID, _ := strconv.Atoi(r.FormValue("recipe-id"))

	fmt.Printf("Deleting recipe %d\n", recipeID)

	w.WriteHeader(http.StatusOK)
	return nil
}

func (c *Config) CookAddRecipeHandler(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	lithiumAmount, _ := strconv.Atoi(r.FormValue("new-recipe-lithium"))
	acetoneAmount, _ := strconv.Atoi(r.FormValue("new-recipe-acetone"))
	sulfuricAmount, _ := strconv.Atoi(r.FormValue("new-recipe-sulfuric"))
	recipeCount, _ := strconv.Atoi(r.FormValue("recipe-count"))

	recipe := Recipe{
		ID: recipeCount + 1,
		Ingredients: []Ingredient{
			{Name: "Lithium", Amount: lithiumAmount},
			{Name: "Acetone", Amount: acetoneAmount},
			{Name: "Sulfuric Acid", Amount: sulfuricAmount},
		},
		IsSelected: false,
	}

	return c.tpl.ExecuteTemplate(w, "cookRecipe", recipe)
}

func (c *Config) StorageAddCardHandler(w http.ResponseWriter, r *http.Request) error {
	cookie, _ := r.Cookie("session_token")
	sessionID := cookie.Value
	sessionUUID, _ := uuid.Parse(sessionID)

	session, _ := c.SessionStore.GetSession(sessionUUID)
	groupName := session.LoginID
	storage, err := c.StorageToolStore.CreateStorage(groupName)
	if err != nil {
		log.Fatalln("failed to create new storage", err)
	}
	data := map[string]interface{}{
		"ID":           storage.ID,
		"PropertyName": storage.StorageName,
		"SlabsClear":   storage.ClearSlabQty.Int32,
		"BlocksClear":  storage.ClearBlockQty.Int32,
		"SlabsCloudy":  storage.CloudySlabQty.Int32,
		"BlocksCloudy": storage.CloudyBlockQty.Int32,
	}

	buf := bytes.Buffer{}
	c.tpl.ExecuteTemplate(&buf, "storagePropertyCard", data)

	msgData := map[string]string{
		"Action":     "Add",
		"HX-Trigger": buf.String(),
	}

	conns := c.conns[groupName]
	for _, conn := range conns {
		if sessionID != conn.sessionID {
			conn.conn.WriteJSON(msgData)
		}
	}

	return c.tpl.ExecuteTemplate(w, "storagePropertyCard", data)
}

func (c *Config) StorageDeleteCardHandler(w http.ResponseWriter, r *http.Request) error {
	cookie, _ := r.Cookie("session_token")
	sessionID := cookie.Value
	sessionUUID, _ := uuid.Parse(sessionID)
	session, _ := c.SessionStore.GetSession(sessionUUID)
	groupName := session.LoginID

	storageID, _ := strconv.Atoi(r.FormValue("storage-card-id"))

	err := c.StorageToolStore.DeleteStorage(groupName, storageID)
	if err != nil {
		fmt.Println("error deleting storage", err)
	}

	msgData := map[string]string{
		"Action":     "Delete",
		"HX-Trigger": fmt.Sprintf("storage-card-%d", storageID),
	}

	conns := c.conns[groupName]
	for _, conn := range conns {
		if sessionID != conn.sessionID {
			conn.conn.WriteJSON(msgData)
		}
	}

	w.WriteHeader(http.StatusOK)
	return nil
}

func (c *Config) LoginHandler(w http.ResponseWriter, r *http.Request) error {
	err := r.ParseForm()
	if err != nil {
		log.Println("error parsing form:", err)
		return err
	}

	groupName := r.FormValue("group-name")
	password := r.FormValue("password")

	isLoggedIn := c.SessionStore.Login(groupName, password)
	if !isLoggedIn {
		data := map[string]interface{}{
			"IsAuthorized": false,
			"WsUrl":        c.WsUrl,
		}
		return c.tpl.ExecuteTemplate(w, "storageTool", data)
	}

	// c.StorageToolStore.Connect(groupName)

	// create sessions
	sessionID := uuid.New()
	expiresAt := time.Now().Add(time.Hour * 24)
	createdAt := time.Now()

	err = c.SessionStore.CreateSession(sessionID, groupName, createdAt, expiresAt)
	if err != nil {
		log.Fatal(err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionID.String(),
		Expires:  expiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	storages, err := c.StorageToolStore.GetGroupStorage(groupName)
	if err != nil {
		log.Fatal("failed to get storages", err)
	}

	data := map[string]interface{}{}

	propertyInventories := []map[string]interface{}{}

	for _, storage := range storages {
		storageData := map[string]interface{}{
			"ID":           storage.ID,
			"PropertyName": storage.StorageName,
			"SlabsClear":   storage.ClearSlabQty.Int32,
			"BlocksClear":  storage.ClearBlockQty.Int32,
			"SlabsCloudy":  storage.CloudySlabQty.Int32,
			"BlocksCloudy": storage.CloudyBlockQty.Int32,
		}
		propertyInventories = append(propertyInventories, storageData)
	}
	data["PropertyInventories"] = propertyInventories
	data["IsAuthorized"] = true
	data["WsUrl"] = c.WsUrl
	return c.tpl.ExecuteTemplate(w, "storageTool", data)
}

func (c *Config) LogoutHandler(w http.ResponseWriter, r *http.Request) error {
	cookie, _ := r.Cookie("session_token")
	sessionID := cookie.Value

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionID,
		Expires:  time.Unix(0, 0),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	data := map[string]interface{}{
		"IsAuthorized": false,
	}
	return c.tpl.ExecuteTemplate(w, "index", data)
}

func (c *Config) GetRegisterHandler(w http.ResponseWriter, r *http.Request) error {
	return c.tpl.ExecuteTemplate(w, "registerUser", nil)
}

func (c *Config) PostRegisterHandler(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		log.Fatal("error parsing form", err)
	}

	groupName := r.FormValue("group-name")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm-password")

	type RegisterUserData struct {
		Errors    map[string]string
		GroupName string // preserve entered data if needed
	}

	errors := make(map[string]string)
	if groupName == "" {
		errors["groupName"] = "Group name is required."
	}
	if password == "" {
		errors["password"] = "Password is required."
	} else if password != confirmPassword {
		errors["confirmPassword"] = "Passwords do not match."
	}

	if len(errors) > 0 {
		data := RegisterUserData{
			Errors:    errors,
			GroupName: groupName, // Optional: pre-fill the form
		}
		return c.tpl.ExecuteTemplate(w, "registerUser", data)
	}

	pwHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	_, err = c.SessionStore.CreateUser(groupName, string(pwHash))
	if err != nil {
		log.Fatal(err)
	}

	return c.tpl.ExecuteTemplate(w, "storageTool", nil)
}

func (c *Config) StorageToolHandleWs(w http.ResponseWriter, r *http.Request) error {
	cookie, _ := r.Cookie("session_token")
	sessionID := cookie.Value
	sessionUUID, _ := uuid.Parse(sessionID)
	session, _ := c.SessionStore.GetSession(sessionUUID)
	groupName := session.LoginID

	upgrader := websocket.Upgrader{
		ReadBufferSize:  0,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	go readPump(conn)
	cn := Connection{
		conn:      conn,
		sessionID: sessionID,
	}

	// for _, cnn := range c.conns[groupName] {
	// 	if cnn.sessionID == sessionID {
	// 		cnn.conn.Close()
	// 	}
	// }

	c.conns[groupName] = append(c.conns[groupName], cn)

	return nil
}

func (c *Config) StorageToolUpdateHandler(w http.ResponseWriter, r *http.Request) error {
	cookie, _ := r.Cookie("session_token")
	sessionID := cookie.Value
	sessionUUID, _ := uuid.Parse(sessionID)
	session, _ := c.SessionStore.GetSession(sessionUUID)
	groupName := session.LoginID

	err := r.ParseForm()
	if err != nil {
		fmt.Println("error parsing form: ", err)
		return err
	}

	storageID, _ := strconv.Atoi(r.FormValue("storage-card-id"))
	storageName := r.FormValue("storage-name")
	clearSlabsQty, _ := strconv.Atoi(r.FormValue("clear-slab-count"))
	clearBlocksQty, _ := strconv.Atoi(r.FormValue("clear-block-count"))
	cloudySlabsQty, _ := strconv.Atoi(r.FormValue("cloudy-slab-count"))
	cloudyBlocksQty, _ := strconv.Atoi(r.FormValue("cloudy-block-count"))

	data := map[string]interface{}{
		"ID":           storageID,
		"PropertyName": storageName,
		"SlabsClear":   clearSlabsQty,
		"BlocksClear":  clearBlocksQty,
		"SlabsCloudy":  cloudySlabsQty,
		"BlocksCloudy": cloudyBlocksQty,
	}

	err = c.StorageToolStore.UpdateStorage(groupName, storageID, storageName, clearSlabsQty, clearBlocksQty, cloudySlabsQty, cloudyBlocksQty)
	if err != nil {
		fmt.Println("error updating storage")
		return err
	}
	buf := &bytes.Buffer{}
	c.tpl.ExecuteTemplate(buf, "storagePropertyCard", data)

	broadcastTemplate(sessionID, c.conns[groupName], buf.Bytes(), websocket.TextMessage)
	return nil
}
