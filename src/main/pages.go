package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mattn/go-sqlite3"
)

var DB_DRIVER string

func init() {
	sql.Register(DB_DRIVER, &sqlite3.SQLiteDriver{})
}

// Não consegui encontrar uma maneira de fazer um forloop simples no template do golang, logo vou exportar como estrutra de dados.
func createPaginationSlice(end int) []int {
	var slice []int
	for i := 0; i < end; i++ {
		slice = append(slice, i+1)
	}
	return slice
}

func (our Developers) checkForDuplicates(other Developer) bool {
	for _, this := range our {
		if this.Name == other.Name {
			return true
		}
	}
	return false
}

/* Index:
Página principal.
*/
func Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/developers", http.StatusFound)
}

/* Developers:
Mostra página que lista (paginado) todos os desenvolvedores.
*/
func showAllDevelopers(w http.ResponseWriter, r *http.Request) {
	// Vamos inicializar uma sessão para que possamos visualizar o "shopping cart" de desenvolvedores
	session, err := store.Get(r, "shopping-cart")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	cookie := session.Values["Cart"]
	var cart = &Cart{}
	var ok bool
	// Abaixo é uma verificação se o que recebemos da sessão é de fato uma estrutura do tipo Cart.
	if cart, ok = cookie.(*Cart); !ok {
		cart = &Cart{Selected_Developers: make([]Developer, 0), Total_Price: 0.0}
	}

	// Valores atuais da sessão
	selected_developers := cart.Selected_Developers
	total_price := cart.Total_Price
	total_price_s := fmt.Sprintf("%.2f", total_price)
	num_devs := len(selected_developers)

	var cart_exists bool
	if num_devs > 0 {
		cart_exists = true
	} else {
		cart_exists = false
	}

	// Acesso ao banco de dados com as informações de um repositório do GitHub.
	database, err := sql.Open(DB_DRIVER, "../tools/db/google_gson.db")
	if err != nil {
		fmt.Println("Failed to open developers database")
	}

	var offset int
	var current_page int
	page := r.URL.Query().Get("p")
	if page != "" {
		i, err := strconv.Atoi(page)
		current_page = i
		if err != nil {
			log.Fatal(err)
		}
		offset = 10 * (i - 1)
	} else {
		offset = 0
		current_page = 1
	}
	session.Values["Last_page"] = current_page
	session.Save(r, w)

	// Vamos calcular também o total de páginas necessárias para mostrar todos os resultados.
	var count int
	rows, err := database.Query("SELECT Count(*) FROM developers")
	if err != nil {
		log.Fatal(err)
	}

	rows.Next()
	err = rows.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	n_pages := int(math.Ceil(float64(count) / 10.0))
	sql_query := fmt.Sprintf("SELECT * FROM developers ORDER BY Lower(Name) ASC LIMIT %d,10", offset)
	rows, err = database.Query(sql_query)
	if err != nil {
		log.Fatal(err)
	}

	devs := make(Developers, 0)
	for rows.Next() {
		var empName sql.NullString
		var empFollowers sql.NullInt64
		var empStars sql.NullInt64
		var empCommits sql.NullInt64
		var empN_repos sql.NullInt64
		var empAvatar sql.NullString

		if err := rows.Scan(&empName, &empFollowers, &empStars, &empCommits, &empN_repos, &empAvatar); err != nil {
			log.Fatal(err)
		}
		var price float64
		price = (float64(empFollowers.Int64)/10.0 + float64(empStars.Int64)/5.0 + float64(empCommits.Int64)/2.0 + float64(empN_repos.Int64)/5.0) + 100
		price_s := fmt.Sprintf("%.2f", price)
		devs = append(devs, Developer{Name: empName.String, Followers: empFollowers.Int64, Stars: empStars.Int64, Commits: empCommits.Int64, N_repos: empN_repos.Int64, Avatar: empAvatar.String, Price: price_s})
	}

	// Como t.Execute espera apenas uma estrutura de dados como parâmetro, precisamos criar uma customizada que junta n_pages com devs.
	data := struct {
		Current_page int
		N_pages      int
		Pagination   []int
		Developers   Developers
		Cart_items   int
		Cart_total   string
		Cart_exists  bool
	}{
		Current_page: current_page,
		N_pages:      n_pages,
		Pagination:   createPaginationSlice(n_pages),
		Developers:   devs,
		Cart_items:   num_devs,
		Cart_total:   total_price_s,
		Cart_exists:  cart_exists,
	}

	t, err := template.ParseFiles("templates/developers.html")
	if err != nil {
		panic(err)
	} else {
		t.Execute(w, data)
	}

}

/* showDeveloper:
Mostra página de desenvolvedor individual
*/
func showDeveloper(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	devId := vars["devId"]
	fmt.Fprintln(w, "Developer #", devId)
}

/* updateCart
Página escondida do usuário (acessível apenas por uso de FORM), que serve para introduzir novo desenvolvedor ao carrinho de compras.
*/
func updateCart(w http.ResponseWriter, r *http.Request) {
	// Vamos inicializar uma sessão para que possamos atualizar o "shopping cart" de desenvolvedores
	session, err := store.Get(r, "shopping-cart")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	cookie := session.Values["Cart"]
	var cart = &Cart{}
	var ok bool
	// Abaixo é uma verificação se o que recebemos da sessão é de fato uma estrutura do tipo Cart.
	if cart, ok = cookie.(*Cart); !ok {
		cart = &Cart{Selected_Developers: make([]Developer, 0), Total_Price: 0.0}
	}

	// Valores antigos da sessão
	selected_developers := cart.Selected_Developers
	old_price := cart.Total_Price

	// Valores novos
	r.ParseForm()
	dev_name := r.Form["addToCart"][0]
	dev_price := r.Form["price"][0]

	// dev_price é uma string, necessário conversão.
	dev_price_f, err := strconv.ParseFloat(dev_price, 64)
	if err != nil {
		log.Fatal(err)
	}

	// Verificação de duplicados: gera um booleando afirmando se já existe o mesmo desenvolvedor ou não na sessão.
	other_dev := Developer{Name: dev_name, Price: dev_price}
	alreadyExists := selected_developers.checkForDuplicates(other_dev)

	url := fmt.Sprintf("/developers?p=%d", session.Values["Last_page"])
	// Atualização da sessão
	if alreadyExists {
		// Redirect de volta ao shop na página que estava.
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		session.Values["Cart"] = Cart{
			Selected_Developers: append(selected_developers, other_dev),
			Total_Price:         old_price + dev_price_f,
		}

		session.Save(r, w)
		http.Redirect(w, r, url, http.StatusFound)
	}

}

/* showCart
Página que lista o carrinho de compras do usuário.
*/
func showCart(w http.ResponseWriter, r *http.Request) {
	// Vamos inicializar uma sessão para que possamos visualizar o "shopping cart" de desenvolvedores
	session, err := store.Get(r, "shopping-cart")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	cookie := session.Values["Cart"]
	var cart = &Cart{}
	var ok bool
	// Abaixo é uma verificação se o que recebemos da sessão é de fato uma estrutura do tipo Cart.
	if cart, ok = cookie.(*Cart); !ok {
		cart = &Cart{Selected_Developers: make([]Developer, 0), Total_Price: 0.0}
	}

	// Invés de enviar a estrutura cart para o template, vamos mandar uma versão processada.
	// Por exemplo, é melhor colocar o preço com duas casas decimais. Na lista de desenvolvedores
	// podemos também fazer com que duplicados tenham um campo de duplicidade.
	price_s := fmt.Sprintf("%.2f", cart.Total_Price)
	data := struct {
		Total_Price string
		Developers  Developers
	}{
		Total_Price: price_s,
		Developers:  cart.Selected_Developers,
	}

	t, err := template.ParseFiles("templates/cart.html")
	if err != nil {
		panic(err)
	} else {
		t.Execute(w, data)
	}
}

/*removeFromCart
Página escondida do usuário (acessível apenas por uso de FORM), que serve para remover desenvolvedor do carrinho de compras.
*/
func removeFromCart(w http.ResponseWriter, r *http.Request) {
	// Vamos inicializar uma sessão para que possamos atualizar o "shopping cart" de desenvolvedores
	session, err := store.Get(r, "shopping-cart")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Remoção do desenvolvedor escolhido do carrinho de compras.
	r.ParseForm()
	dev_name := r.Form["removeFromCart"][0]
	dev_price := r.Form["price"][0]

	// dev_price é uma string, necessário conversão.
	dev_price_f, err := strconv.ParseFloat(dev_price, 64)
	if err != nil {
		log.Fatal(err)
	}

	cookie := session.Values["Cart"]
	var cart = &Cart{}
	var ok bool
	// Abaixo é uma verificação se o que recebemos da sessão é de fato uma estrutura do tipo Cart.
	if cart, ok = cookie.(*Cart); !ok {
		cart = &Cart{Selected_Developers: make([]Developer, 0), Total_Price: 0.0}
	}

	cart_devs := cart.Selected_Developers
	cart_price := cart.Total_Price

	for i, dev := range cart_devs {
		if dev.Name == dev_name {
			cart_devs = append(cart_devs[:i], cart_devs[i+1:]...)
			cart_price = cart_price - dev_price_f
			break
		}
	}

	// Agora que removemos o desenvolvedor não requerido, podemos atualizar a sessão.
	// Também  necessário verificar agora quantos desenvolvedores temos na nossa lista.
	// Se não houver nenhum, redireciona à loja. Zerar também o preço total (por causa de arrendondamento
	// mesmo se removermos todos os desenvolvedores, Total_Price não será exatamente zero.
	n_devs := len(cart_devs)
	if n_devs > 0 {
		session.Values["Cart"] = Cart{
			Selected_Developers: cart_devs,
			Total_Price:         cart_price,
		}
		session.Save(r, w)
		http.Redirect(w, r, "/check-out", http.StatusFound)
	} else {
		session.Values["Cart"] = Cart{
			Selected_Developers: cart_devs,
			Total_Price:         0.0,
		}
		session.Save(r, w)
		http.Redirect(w, r, "/developers", http.StatusFound)
	}
}

/*confirmPurchase
Página escondida do usuário (acessível apenas por uso de FORM), que confirma ao usuário que ele de fato comprou o serviço
dos desenvolvedores queridos.
*/
func confirmPurchase(w http.ResponseWriter, r *http.Request) {
	// Vamos inicializar uma sessão para que possamos visualizar o "shopping cart" de desenvolvedores
	session, err := store.Get(r, "shopping-cart")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	cookie := session.Values["Cart"]
	var cart = &Cart{}
	var ok bool
	// Abaixo é uma verificação se o que recebemos da sessão é de fato uma estrutura do tipo Cart.
	if cart, ok = cookie.(*Cart); !ok {
		cart = &Cart{Selected_Developers: make([]Developer, 0), Total_Price: 0.0}
	}

	cart_devs := cart.Selected_Developers

	// Vamos colher mais informações  desses desenvolvedores
	// Primeiro vamos criar o acesso ao banco de dados
	database, err := sql.Open(DB_DRIVER, "../tools/db/google_gson.db")
	if err != nil {
		fmt.Println("Failed to open developers database")
	}

	sql_query := "SELECT * FROM developers WHERE"
	num_devs := len(cart_devs)
	for i, dev := range cart_devs {
		sql_query = fmt.Sprintf("%s Name='%s'", sql_query, dev.Name)
		if num_devs > i+1 {
			sql_query = fmt.Sprintf("%s OR", sql_query)
		}
	}

	fmt.Println(sql_query)
	rows, err := database.Query(sql_query)
	if err != nil {
		log.Fatal(err)
	}

	devs := make(Developers, 0)
	for rows.Next() {
		var empName sql.NullString
		var empFollowers sql.NullInt64
		var empStars sql.NullInt64
		var empCommits sql.NullInt64
		var empN_repos sql.NullInt64
		var empAvatar sql.NullString

		if err := rows.Scan(&empName, &empFollowers, &empStars, &empCommits, &empN_repos, &empAvatar); err != nil {
			log.Fatal(err)
		}

		devs = append(devs, Developer{Name: empName.String, Followers: empFollowers.Int64, Stars: empStars.Int64, Commits: empCommits.Int64, N_repos: empN_repos.Int64, Avatar: empAvatar.String, Price: ""})
	}

	// Agora que confiramos a compra, precisamos destruir a sessão.
	session.Values["Cart"] = nil
	session.Save(r, w)

	t, err := template.ParseFiles("templates/confirm.html")
	if err != nil {
		panic(err)
	} else {
		t.Execute(w, devs)
	}
}
