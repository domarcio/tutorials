# Refatoração e testes

Cada línguaguem de programação tem uma forma única de implementar testes e com `Go` não é diferente.

Deixaremos todo o debate que existe acerca sobre testes e focaremos única e exclusivamente em como trabalhar com eles usando Go.

O **objetivo** desse pequeno tutorial é te dar uma base sobre *(1)* o que são testes em Go e como estruturá-los, como executá-los e por fim *(2)* fazer uma pequena refatoração em um código aplicando os conceitos de [**TDD**](https://pt.wikipedia.org/wiki/Test-driven_development).

## Sobre testes

Go já tem uma abordagem intrínseca para se trabalhar com testes, um simples comando como `go test ./...` é o suficiente para execução de todos os testes da aplicação.

Há dois modos diferentes para execução dos testes.

O primeiro modelo é quando executamos o comando `go test` sem nenhum argumento. Dessa forma, `go test` compilará o código-fonte do pacote e os testes encontrados no diretório onde o comando é executado, em seguida teremos o resultado desse binário compilado.

```bash
$ go test
PASS
ok      your-project 0.176s
```

O segundo modelo é quando executamos o teste de um pacote em específico com a seguinte sintaxe: `go test packagename` (ou `go test ./...` como visto acima para executar todos os testes em todos os pacotes presentes recursivamente). Assim como na primeira abordagem, o comando compilará o código-fonte e os testes de todo(s) o(s) pacote(s) que foi(foram) encontrado(s) para em seguida termos o resultado desse binário compilado.

```bash
$ go test ./pkg1/pkg2/wallet/cashout
ok      your-project/pkg1/pkg2/wallet/cashout       0.009s
```

Toda saída dos testes executados retornarão sempre dois resultados possíveis: `ok` ou `FAIL` para cada pacote verificado.

```bash
$ go test ./...
ok   archive/tar   0.011s
FAIL archive/zip   0.022s
ok   compress/gzip 0.033s
```

Alguns flags comumente usadas são `-v`, `-bench` e `-race` junto ao comando `go test`. Para mais detalhes veja a documentação oficial: `go help test`.

**Uma nota sobre Cache:** Após a execução dos testes o Go sempre fará um cache do binário para evitar execução desnecessária. Sempre que você executar um teste que está em cache a seguinte assinatura `(cached)` é exibida no resultado final.

```bash
$ go test ./pkg1/pkg2/wallet/cashout
ok      your-project/pkg1/pkg2/wallet/cashout       0.009s       (cached)
```

### Criando nossos testes

O comando `go test` busca sempre pelos arquivos nomeados com `*_test.go` e faz os teste de todas as funções dos tipos *test* e *benchmarks*.

Cada função para ser lida pelo comando precisa ter sua assinatura que deve sempre estar sempre presente no ínicio do nome na função:

1. Funções de teste: `func TestWordReverse(...)`.
2. Funções de desempenho: `func BenchmarkWordReverse(...)`.

Todo arquivo de teste precisa ter (obrigatoriamente) o módulo `testing` declarado no seus *imports*:

```go
package foo

import "testing"
```

Segue alguns exemplos (abaixo) de testes para a nossa função `WordReverse`:

```go
package examples

func WordReverse(s string) string {
    var (
        runes  = []rune(s)
        strLen = len(runes)
    )

    for i, j := 0, strLen-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }

    return string(runes)
}
```

#### Funções de Test

Como já dito anteriormente, toda função de teste deve começar com o prefixo `Test` e seguir (sufixo) com o nome da função em si. Exemplo:

```go
package examples

func TestWordReverse(t *testing.T) {
    // ...
}
```

**Importante** reparar no parâmetro `t *testing.T` que a função recebe. Esse argumento nos fornece métodos onde informamos falhas nos testes, logging e assim por diante, para mais detalhes [veja a documentação](https://golang.org/pkg/testing/#T).

Por fim um outro exemplo prático de como se testar uma função.

```go
package examples

func TestWordReverse1(t *testing.T) {
    var (
        expeced = "dlroW olleH"
        result  = WordReverse("Hello World")
    )

    if result != expeced {
        t.Errorf("got %s, expected %s", result, expeced)

        // Veja mais opções na doc
        // t.Error()
        // t.Fail()
        // t.FailNow()
        // ...
    }
}
```

### Funções Benchmark

Assim como nas funções `Test`, as funções de desempenho precisam ter um prefixo `Benchmark` seguido do nome da função. Exemplo:

```go
package examples

func BenchmarkWordReverse(b *testing.B) {
    // ...
}
```

**Importante** reparar no parâmetro `b *testing.B` que a função recebe. Esse argumento nos fornece métodos em sua grande maioria iguais aos da função `t *testing.T` e outros para medir o desempenho da função, para mais detalhes [veja a documentação](https://golang.org/pkg/testing/#B).

Por fim um outro exemplo prático de como se avaliar o desempenho de uma função.

```go
package examples

func BenchmarkWordReverse(b *testing.B) {
    for i := 0; i < b.N; i++ {
        WordReverse("Hello World")
    }
}
```

**Os desempenhos são avaliados apenas quando usamos o argumento `--bench=.` no comando `go test`**. Exemplo: `go test --bench=BenchmarkWordReverse -benchmem`.

## Testes com DynamoDB

Agora vamos a parte mais legal desse pequeno texto: Refatorar uma função que busca um registro no DynamoDB.

Para fins de teste você pode subir um container local com o DynamoDB e criar um item para testarmos. O objetivo desse container é você poder executar o `main.go` antes e depois das alterações, assim garantirá que nada se perdeu.

```bash
~ docker run -d -p 8000:8000 amazon/dynamodb-local
~ aws dynamodb create-table \
    --endpoint-url http://localhost:8000 \
    --region us-east-1 \
    --table-name Transaction \
    --attribute-definitions \
        AttributeName=ID,AttributeType=S \
        AttributeName=CustomerID,AttributeType=S \
    --key-schema \
        AttributeName=ID,KeyType=HASH \
        AttributeName=CustomerID,KeyType=RANGE \
    --provisioned-throughput \
        ReadCapacityUnits=10,WriteCapacityUnits=5
~ aws dynamodb put-item \
    --endpoint-url http://localhost:8000 \
    --region us-east-1 \
    --table-name Transaction \
    --item '{"ID":{"S":"f4302dff-5e50-4d07-b201-8e675a583c2a"},"CustomerID":{"S":"d07b14ae-55dc-4f68-917f-cd9857994b16"},"Amount":{"N":"10.57"}}' \
    --return-consumed-capacity TOTAL \
    --return-item-collection-metrics SIZE
```

> É necessário ter instalado o [aws-cli](https://docs.aws.amazon.com/pt_br/rekognition/latest/dg/setup-awscli-sdk.html) para conseguir rodar os comandos do DynamoDB para o Docker.

Nossa código de exemplo consiste em **buscar uma Transação de acordo com o ID e Customer ID**.

```go
// main.go
package main

import (
    "examples/transaction"
    "log"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
)

var (
    endpoint = "http://localhost:8000"
    region   = "us-east-1"
)

func main() {
    config := &aws.Config{
        CredentialsChainVerboseErrors: aws.Bool(true),
        Endpoint:                      aws.String(endpoint),
        Region:                        aws.String(region),
    }
    sess, err := session.NewSession(config)
    if err != nil {
        log.Fatalln(err)
    }

    db := dynamodb.New(sess)
    repo := transaction.NewRepository(db)

    t, err := repo.GetTransaction("f4302dff-5e50-4d07-b201-8e675a583c2a", "d07b14ae-55dc-4f68-917f-cd9857994b16")
    if err != nil {
        log.Fatalln(err)
    }

    log.Printf("ID: %s, Customer ID: %s, Amount %s", t.ID, t.CustomerID, t.Amount.StringFixed(2))
}

// transaction/repository.go
package transaction

import (
    "errors"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
    ErrTransactionNotFound error = errors.New("transaction does not exists")
)

type Repository struct {
    db *dynamodb.DynamoDB
}

func NewRepository(db *dynamodb.DynamoDB) *Repository {
    return &Repository{db: db}
}

func (r *Repository) GetTransaction(id, customerID string) (*Transanction, error) {
    input := &dynamodb.GetItemInput{
        Key: map[string]*dynamodb.AttributeValue{
            "ID": {
                S: aws.String(id),
            },
            "CustomerID": {
                S: aws.String(customerID),
            },
        },
        TableName: aws.String("Transaction"),
    }

    result, err := r.db.GetItem(input)
    if err != nil {
        return nil, err
    }

    if len(result.Item) <= 0 {
        return nil, ErrTransactionNotFound
    }

    transaction := &Transanction{}
    err = dynamodbattribute.UnmarshalMap(result.Item, transaction)
    if err != nil {
        return nil, err
    }

    return transaction, nil
}
```

O código acima funciona, porém é um pouco díficil para testar por 1 motivo:

1. O repositório faz o uso da implementação do DynamoDB ao invés de usar [suas próprias interfaces](https://docs.aws.amazon.com/sdk-for-go/api/service/dynamodb/dynamodbiface/).

Claro que no "mundo real" é bem comum encontrar esse tipo de situação e não há nada de mal nisso, porém teríamos que subir uma instância do Docker para garantir que nosso teste funcione.

Para driblar esse tipo de situação nós faremos um pequeno refactor e passar a usar as `interfaces` do DynamoDB.

```go
// transaction/repository.go
package transaction

import (
    "errors"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
    ErrTransactionNotFound error = errors.New("transaction does not exists")
)

type Repository struct {
    // Passamos a usar interface do Dynamo
    db dynamodbiface.DynamoDBAPI
}

func NewRepository(db dynamodbiface.DynamoDBAPI) *Repository {
    return &Repository{db: db}
}

func (r *Repository) GetTransaction(id, customerID string) (*Transanction, error) {
    // A implementação continua a mesma
}
```

Se você for no terminal e executar o `main.go` verá que o funcionamento continua o mesmo.

Dessa forma ficou mais simples para se testar, então vamos aos testes:

```go
// transaction/repository_test.go
package transaction

import (
    "errors"
    "reflect"
    "testing"

    "github.com/shopspring/decimal"
)

func TestGetTransaction(t *testing.T) {
    cases := []struct {
        name                string
        transactionID       string
        customerID          string
        expectedTransaction *Transanction
        expectedError       error
    }{
        {
            name:                "an aws error",
            transactionID:       "foo",
            customerID:          "bar",
            expectedTransaction: nil,
            expectedError:       errors.New("aws error"),
        },
        {
            name:                "transaction not found",
            transactionID:       "foo",
            customerID:          "bar",
            expectedTransaction: nil,
            expectedError:       ErrTransactionNotFound,
        },
        {
            name:                "successful",
            transactionID:       "foo",
            customerID:          "bar",
            expectedTransaction: &Transanction{ID: "foo", CustomerID: "bar", Amount: Currency{decimal.New(10, 50)}},
            expectedError:       nil,
        },
    }

    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            repo := NewRepository(nil)
            transaction, err := repo.GetTransaction(c.transactionID, c.customerID)

            if c.expectedError != nil && !errors.Is(c.expectedError, err) {
                t.Errorf("expected error %s, got %s", c.expectedError.Error(), err.Error())
            } else if c.expectedTransaction != nil && !reflect.DeepEqual(c.expectedTransaction, transaction) {
                t.Errorf("expected transaction %+v, got %+v", c.expectedTransaction, transaction)
            }
        })
    }
}
```

Alguns cenários prontos para testarmos nossa função `GetTransaction`. O que precisa ser feito agora é informar ao `repo := NewRepository(nil)` qual implementação usaremos para suprir as necessidades do `db dynamodbiface.DynamoDBAPI`.

Como dito anteriormente, o bom de usarmos `interface` nos parâmetros das funções é que conseguimos criar *mocks* para nossos testes e assim garantiremos que o contrato é respeitado, o fluxo está coerente e saberemos exatamente quais são os dados de entrada/saída da nossa função.

Vamos criar nosso *mock* para o DynamoDB e informá-lo como argumento na variável `repo := NewRepository(nil)`.

```go
// transaction/mock/dynamodb.go
package mock

import (
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type DynamoDBMock struct {
    dynamodbiface.DynamoDBAPI

    GetItemMock func(*dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
}

func (m *DynamoDBMock) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
    return m.GetItemMock(input)
}

```

E agora alterar nosso teste para receber os *mocks*.

```go
// transaction/repository_test.go
t.Run(c.name, func(t *testing.T) {
    mock := &mock.DynamoDBMock{
        GetItemMock: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
            return nil, nil
        },
    }
    repo := NewRepository(mock)
    // ...
})
```

Veja que chamamos uma função `GetItemMock` dentro da nossa struct `mock.DynamoDBMock`, ela é responsável por adicionar saídas (`(*dynamodb.GetItemOutput, error)`) de acordo com cada cenário do teste que precisamos.

Para ficar um pouco mais dinâmico nós adicionaremos um *mock case* em cada um dos nossos cenários de testes na variável `cases`:

```go
// transaction/repository_test.go
func TestGetTransaction(t *testing.T) {
    cases := []struct {
        name                    string
        transactionID           string
        customerID              string
        expectedTransaction     *Transanction
        expectedError           error
        mockDynamoDBGetItemMock func() (*dynamodb.GetItemOutput, error) // Adicionado
    }{
        {
            name:                "an aws error",
            transactionID:       "foo",
            customerID:          "bar",
            expectedTransaction: nil,
            expectedError:       errors.New("aws error"),
            mockDynamoDBGetItemMock: func() (*dynamodb.GetItemOutput, error) {
                return nil, errors.New("aws error")
            },
        },
        {
            name:                "transaction not found",
            transactionID:       "foo",
            customerID:          "bar",
            expectedTransaction: nil,
            expectedError:       ErrTransactionNotFound,
            mockDynamoDBGetItemMock: func() (*dynamodb.GetItemOutput, error) {
                out := &dynamodb.GetItemOutput{
                    Item: map[string]*dynamodb.AttributeValue{},
                }
                return out, nil
            },
        },
        {
            name:                "successful",
            transactionID:       "foo",
            customerID:          "bar",
            expectedTransaction: &Transanction{ID: "foo", CustomerID: "bar", Amount: Currency{decimal.NewFromFloat(10.5)}},
            expectedError:       nil,
            mockDynamoDBGetItemMock: func() (*dynamodb.GetItemOutput, error) {
                out := &dynamodb.GetItemOutput{
                    Item: map[string]*dynamodb.AttributeValue{
                        "ID":         {S: aws.String("foo")},
                        "CustomerID": {S: aws.String("bar")},
                        "Amount":     {N: aws.String("10.5")},
                    },
                }
                return out, nil
            },
        },
    }

    for _, c := range cases {
        t.Run(c.name, func(t *testing.T) {
            mock := &mock.DynamoDBMock{
                GetItemMock: func(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
                    return c.mockDynamoDBGetItemMock() // Adicionado
                },
            }
            repo := NewRepository(mock)
            transaction, err := repo.GetTransaction(c.transactionID, c.customerID)

            if c.expectedError != nil && c.expectedError.Error() != err.Error() {
                t.Errorf("expected error %s, got %s", c.expectedError.Error(), err.Error())
            } else if c.expectedTransaction != nil && !reflect.DeepEqual(c.expectedTransaction, transaction) {
                t.Errorf("expected transaction %+v, got %+v", c.expectedTransaction, transaction)
            }
        })
    }
}
```

Essa é a última adapção para nosso teste funcionar com os *mocks*. Você pode executar essa função de teste sozinha e ver o resultado de cada item.

```bash
~ pwd
/home/username/projects/your-project

~ go clean -testcache && go test -v -run ^TestGetTransaction$ ./transaction
=== RUN   TestGetTransaction
=== RUN   TestGetTransaction/an_aws_error
=== RUN   TestGetTransaction/transaction_not_found
=== RUN   TestGetTransaction/successful
--- PASS: TestGetTransaction (0.00s)
    --- PASS: TestGetTransaction/an_aws_error (0.00s)
    --- PASS: TestGetTransaction/transaction_not_found (0.00s)
    --- PASS: TestGetTransaction/successful (0.00s)
PASS
ok      your-project/transaction       0.003s
```

**Tudo funcionando perfeitamente!**

Para finalziar deixo a dica de usar as libs [stretchr/testify](https://github.com/stretchr/testify) e [vektra/mockery](https://github.com/vektra/mockery) para trabalhar com asserts e mocks, respectivamente.
