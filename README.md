# golang-learning
## Um projeto desafio introduzido pela empresa VTEX.
---
### ABOUT

O desafio se reduz a montar um servidor de páginas dinâmicas, feito em golang, linguagem recentemente desenvolvida pela Google.

Esse servidor, baseado em RESTful princípios, seria o host de um website de contratação de programadores. O "site" tem apenas três páginas distintas: a loja; o carrinho de compras e a confirmaçao da contratação desses programadores.

As informações sobre todos os programadores deveria ser buscada a partir do próprio GitHub, a partir da lista de contribuintes de um projeto qualquer.

O desafio tinha que ser completado em até dois dias, com diversas condições sendo satisfeitas, baseado no que o autor priorizava.

### O QUE FOI FEITO DE FATO

* Considero que foi construído com razoável sucesso um servidor que de fato hospeda o site descrito acima.
* A obtenção dos dados do GitHub foi feita via o desenvolvimento de uma ferramento em python, o populate_db.py. Notou-se a partir desse exercício que há muitas falhas de usabilidade no API do GitHub.
* A persistência dos dados do site, necessária para a implementação de um carrinho de compras, foi feita via sessões e cookies.
* Testes automatizados foram feitos no próprio Go.

### DESCRIÇÃO DA ESTRUTURA DE ARQUIVOS
bin: contém o executável do servidor, "main", como também arquivos estáticos e templates usados pelo servidor.
pkg: gerado automáticamente pelo comando "go get"
src: contém todas as bibliotecas externadas usadas, o próprio código .go do servidor e também testes.
tools: contém as ferramentas auxiliares feito em python como também os bancos de dados que elas geraram.

### CURIOSIDADE
O desafio permitiu que o autor escolhesse sua própria maneira de calcular o preço de cada programador.
O meu é:

100 + 1/10*número_de_followers + 1/5*número_de_estrelas + 1/5*número_de_repos + 1/2*número_de_commits

A idéia é que seguidores denota popularidade, não necessáriamente capacidade, logo um peso pequeno. O número de repositórios é proporcional ao nível de atividade do programador, já o número de estrelas provavelmente é proporcional ao quão creativo são seus projetos. Como ambos os campos denota capacidade cognitiva, dei maior pesos a eles. Por final, dei maior peso às contribuições ao projeto de interesse, já que isso de certa maneira mostra que ele entende daquele assunto/projeto específico. Como isso pode ser de grande importância para a pessoa que o contrata, achei razoável colocar o maior peso a esse item.

### PRÉ-REQUISITOS
go version go1.3

### BUILD
A partir do nível topo do projeto faça:

    cd src/main
    go get
    go install main

### RUN
A partir do nível topo do projeto rode:

    bin/main

### TESTING
A partir do nível topo do projeto faça:

    cd src/tests
    go test --v

