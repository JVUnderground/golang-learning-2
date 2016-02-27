# -*- coding: utf-8 -*-
#
# populate_db(owner, project): returns database with project contributors 

import requests
import sys
import re
import os

from requests.packages.urllib3.exceptions import InsecurePlatformWarning, SNIMissingWarning

requests.packages.urllib3.disable_warnings(InsecurePlatformWarning)
requests.packages.urllib3.disable_warnings(SNIMissingWarning)

# Developer class - see developer.py
from developer import Developer

# Script configuration variables:
# Infelizmente o GitHub API limita o usuário anônimo a apenas 60 'requests'por dia. Assim, é
# necessário o uso de login, senha de algum usuário, o que aumenta a taxa para 5000 por dia.
# Cuidado para retirar essas informações ao compartilhar essa ferramenta com outros!
#
# gh_user - usuário do GitHub
# gh_pass - senha do usuário

gh_user = ""
gh_pass = ""


# getNumPages(header): Retorna um dicionário {n_pages, last_page}
# onde n_pages é o número de páginas que o request contém (máximo de 30 resultados por página)
# e last_page é o URL da última página. A entrada header é o cabecalho de uma resposta HTTP do 
# GitHub API.
#
# A função usa uma variável global auxiliar, last_url, a fim de evitar a repetição muito grande
# da compilação de uma expressão regular.
last_url = re.compile('\s*<(.*page=(.*))>;\s*rel="last"')
def getNumPages(headers):
    n_pages = 1
    last_page = ''
    if 'link' in headers:
        links = headers['link'].split(',')
        
        for link in links:
            match_bool = last_url.match(link)
            if match_bool:
                n_pages = match_bool.groups()[1]
                last_page = match_bool.groups()[0]
                break

    return {'n_pages': int(n_pages), 'last_page': last_page}

# Verificação simples de argumentos.
if len(sys.argv) != 3:
    print("E1: populate_db expects exactly two (2) inputs, received %d." % (len(sys.argv)-1))
    exit(1)

# Provavelmente necessário a sanitização dos argumentos abaixo. Como eu confio em mim mesmo,
# deixarei assim por enquanto.
owner = sys.argv[1]
repos = sys.argv[2]

db_name = 'db'+ os.sep + owner + '_' + repos + '.db'

request_url = "https://api.github.com/repos/%s/%s/contributors" % (owner, repos)
print("Downloading contributors list at %s." % request_url)
r = requests.get(request_url, auth=(gh_user, gh_pass))

if r.status_code != 200:
    print("E2: repository not found. Status code: %s" % r.status_code)
    exit(2)


# Primeiro vamos preencher nossa variável de desenvolvedores com o resultado da primeira página
# Há alguns problemas no API do GitHub, e.g. não há maneira fácil de obter o número de seguidores.
# Mesma coisa ocorre com o número de estrelas.
#
# O GitHub API ainda trabalha com paginação, então para calcular o número de estrelas é
# necessário fazer um 'request' para cada repositório. Já para o seguidores,
# basta ver a quantidade de páginas - 1, multiplicar por 30 (número máximo por página), e somar
# a quantidade de seguidores na última página.
#
# Como o GitHub tem limite de acesso diário do API (5000), não é recomendado o uso dessa ferramenta
# para repositórios muito populares.
developers = []
resp = r.json()
for dev in resp:
    name = dev['login']
    avatar = dev['avatar_url']
    commits = dev['contributions']

    # Necessário 2 requests para conseguir o número de followers.
    r_fol = requests.get("https://api.github.com/users/%s/followers" % name, auth=(gh_user, gh_pass))
    h_fol = r_fol.headers

    pages_fol = getNumPages(h_fol)
    
    if(pages_fol['n_pages'] > 1):
        r_last_fol = requests.get(pages_fol['last_page'], auth=(gh_user, gh_pass))
        followers = 30*(pages_fol['n_pages']-1) + len(r_last_fol.json())
    else:
        followers = len(r_fol.json())

    # Agora requistando todos os repositórios do usuário para contar o número de estrelas.
    r_repos = requests.get("https://api.github.com/users/%s/repos" % name, auth=(gh_user, gh_pass))
    j_repos = r_repos.json()
    h_repos = r_repos.headers

    stars = 0
    n_repos = 0
    for repos in j_repos:
        stars += repos['watchers']
        n_repos += 1
    
    # Há mais que uma página de resultados?
    pages_repos = getNumPages(h_repos)
    n_pages = pages_repos['n_pages']
    repos_url = pages_repos['last_page']
    
    for page in range(1,n_pages):
        page_num = page + 1 # Já pegamos a primeira página
        repos_url_page = re.sub("page=\d+", "page=%s" % page_num, repos_url)

        r_repos = requests.get(repos_url_page, auth=(gh_user, gh_pass))
        j_repos = r_repos.json()
        # Não precisa mais de headers, pois já sabemos quantas páginas temos que verificar

        for repos in j_repos:
            stars += repos['watchers']
            n_repos += 1

    # Finalmente temos tudo do desenvolvedor individual. Seu nome, seguidores, estrelas
    # contribuições no repositório de interesse e o número de repositórios do qual é dono.
    developers.append(Developer(name, followers, stars, commits, n_repos, avatar))


# Tudo acima apenas foi para a primeira página, se houver mais que uma, é necessário repetir,

headers = r.headers # Lembrando que r é o request original.
pages = getNumPages(headers)
n_pages = pages['n_pages']
contr_url = pages['last_page']
for page in range(1,n_pages):
    page_num = page + 1 # Já pegamos a primeira página
    contr_url_page = re.sub("page=\d+", "page=%s" % page_num, contr_url)
    
    r_contr = requests.get(contr_url_page, auth=(gh_user, gh_pass))
    j_contr = r_contr.json()
    # Não preicsa mais de headers, pois já sabemos quantas páginas temos que verificar.

    for dev in j_contr:
        name = dev['login']
        avatar = dev['avatar_url']
        commits = dev['contributions']

        # Necessário 2 requests para conseguir o número de followers.
        r_fol = requests.get("https://api.github.com/users/%s/followers" % name, auth=(gh_user, gh_pass))
        h_fol = r_fol.headers
    
        pages_fol = getNumPages(h_fol)
       
        if(pages_fol['n_pages'] > 1):
            r_last_fol = requests.get(pages_fol['last_page'], auth=(gh_user, gh_pass))
            followers = 30*(pages_fol['num_pages']-1) + len(r_last_fol.json())
        else:
            followers = len(r_fol.json())
        
        # Agora requistando todos os repositórios do usuário para contar o número de estrelas.
        r_repos = requests.get("https://api.github.com/users/%s/repos" % name, auth=(gh_user, gh_pass))
        j_repos = r_repos.json()
        h_repos = r_repos.headers

        stars = 0
        n_repos = 0
        for repos in j_repos:
            stars += repos['watchers']
            n_repos += 1
        
        # Há mais que uma página de resultados?
        pages_repos = getNumPages(h_repos)
        n_pages = pages_repos['n_pages']
        repos_url = pages_repos['last_page']
        
        for page in range(1,n_pages):
            page_num = page + 1 # Já pegamos a primeira página
            repos_url_page = re.sub("page=\d+", "page=%s" % page_num, repos_url)
    
            r_repos = requests.get(repos_url_page, auth=(gh_user, gh_pass))
            j_repos = r_repos.json()
            # Não precisa mais de headers, pois já sabemos quantas páginas temos que verificar
    
            for repos in j_repos:
                stars += repos['watchers']
                n_repos += 1
    
        # Finalmente temos tudo do desenvolvedor individual. Seu nome, avatar, seguidores, estrelas,
        # contribuições no repositório de interesse e o número de repositórios do qual é dono.
        developers.append(Developer(name, followers, stars, commits, n_repos, avatar))
        
# Temos agora todos os desenvolvedores, contribuidores do repositório de interesse. Precisamos 
# agora armazenar essas informações em um banco de dados.

if os.path.isfile(db_name):
    os.remove(db_name)
    
Developer.init_db(db_name)
for dev in developers:
    dev.insert_to_db()
    
exit(0) # No error.
