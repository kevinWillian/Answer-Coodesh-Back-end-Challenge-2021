 >  This is a challenge by [Coodesh](https://coodesh.com/)

# Answer Coodesh Back end Challenge 2021
## By Kevin Willian Dos Santos

## Introdução:

Esta é a resposta desenvolvida por Kevin Willian Dos Santos ao desafio lançado pela Coodesh chamado: "Back-end Challenge 🏅 2021 - Space Flight News". Nesta resposta temos uma API sincronizada com a Space Flight News.

# Tecnologias:

Nesta resposta usei: **Go Lang | Postgress | Docker** com bibliotecas para GO: **go-chi | robfig-cron | rs-cors |driver-postgres | Gorm**


### Como instalar e usar a resposta:

 * O Docker deve estar devidamente instalado e configurado na maquina usada

 * Configure o arquivo 'conf.json' com as configurações do banco de dados desejado

 * Na Raiz do projeto rode o comando:
    ```

    docker build --tag docker-answer-coodesh .

    ```
 * Em seguida rode a maquina Docker criada expondo a porta usada pela API no host:

    ```

    docker run --publish 3000:3000 docker-answer-coodesh

    ```

 (dica) * leia o arquivo 'docs/api.md' para ver funções API disponiveis

 ## Para subir alterações no git antes execute (na raiz do projeto)


    git update-index --assume-unchanged conf.json

