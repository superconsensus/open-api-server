package sysinit

import (
	"context"
	"fmt"
	"github.com/superconsensus/matrix-sdk-go/v2/xuper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

// todo db初始化

type MongoClient struct {
	*mongo.Client
	*mongo.Database
}

type Count struct {
	//ID        primitive.ObjectID `bson:"_id,omitempty"`
	TxCount   int64  `bson:"tx_count"`   //交易总数
	CoinCount string `bson:"coin_count"` //全网金额
	AccCount  int64  `bson:"acc_count"`  //账户总数
	Accounts  bson.A `bson:"accounts"`   //账户列表
}

// 94db的account表结构
type Account struct {
	Id        string `bson:"_id"`
	Pwd       string `bson:"pwd"`
	Address   string `bson:"address"`
	Timestamp int64  `bson:"timestamp"`
}

var counts *Count

func NewMongoClient(dataSource, database string) (*MongoClient, error) {
	client, err := mongo.NewClient(options.Client().
		ApplyURI(dataSource).
		SetConnectTimeout(10 * time.Second))
	if err != nil {
		return nil, err
	}

	//ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	//defer client.Disconnect(ctx)

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return &MongoClient{client, client.Database(database)}, nil
}

func (m *MongoClient) GetUserAddress(db string, name string) (string, error) {
	userCol := m.Database.Collection(db)
	singleRes := userCol.FindOne(context.TODO(), bson.M{
		"_id": name,
	})
	if singleRes.Err() == nil {
		user := &Account{}
		err := singleRes.Decode(user)
		if err != nil {
			return "", err
		}
		return user.Address, nil
	}
	return "", singleRes.Err()
}

// 测试用
func (m *MongoClient) findAccount() {
	accCol := m.Database.Collection("account")

	if counts == nil {
		counts = new(Count)
	}

	res := accCol.FindOne(nil, bson.M{
		"_id": "123",
	})
	fmt.Println("ressssss", res, res.Err())
	if res != nil {
		a := &Account{}
		res.Decode(a)
		fmt.Println("decode", a)
	}

	cursor, err := accCol.Find(nil, bson.M{})
	if err != nil && err != mongo.ErrNoDocuments {
		return
	}
	if cursor != nil {
		err = cursor.All(nil, &counts.Accounts)
	}
	//过滤key,减小体积
	for i, v := range counts.Accounts {
		counts.Accounts[i] = v.(bson.D).Map()["_id"]
	}

	xclient, _ := xuper.New("120.79.69.94:37102")
	cnt := 0
	for _, account := range counts.Accounts {
		balance, _ := xclient.QueryBalance(account.(string), xuper.WithQueryBcname("xuper"))
		fmt.Println("账号", account, "账号长度", len(account.(string)) == 33, "余额", balance)
		if len(account.(string)) != 33 && len(account.(string)) != 41 {
			continue
		}
		//if balance.Int64() != 0 {
		cnt++
		//}
	}
	fmt.Println("有效地址数量", cnt)
}
