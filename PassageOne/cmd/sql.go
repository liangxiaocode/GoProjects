package cmd

import (
	"PassageOne/internal/sql2struct"
	"log"

	"github.com/spf13/cobra"
)

var (
	username  string
	password  string
	host      string
	charset   string
	dbtype    string
	dbname    string
	tablename string
)

var sqlcmd = &cobra.Command{
	Use:   "sql",
	Short: "sql转换和处理",
	Long:  "sql转换和处理",
	Run:   func(cmd *cobra.Command, args []string) {},
}

var sql2structCmd = &cobra.Command{
	Use:   "struct",
	Short: "sql转struct",
	Long:  "sql转struct",
	Run: func(cmd *cobra.Command, args []string) {
		dbinfo := &sql2struct.DBInfo{
			DBType:   dbtype,
			Host:     host,
			UserName: username,
			PassWord: password,
			Charset:  charset,
		}
		dbModel := sql2struct.NewDBModel(dbinfo)
		err := dbModel.Connect()
		if err != nil {
			log.Fatalf("dbModel.Connect failed,err:%v\n", err)
		}
		columns, err := dbModel.GetColumn(dbname, tablename)
		if err != nil {
			log.Fatalf("dbModel.GetColumns failed,err:%v\n", err)
		}
		tpl := sql2struct.NewStructTemplate()
		tplColumns := tpl.AssemblyColumns(columns)
		err = tpl.Generate(tablename, tplColumns)
		if err != nil {
			log.Fatalf("dbModel.Generate failed,err:%v\n", err)
		}
	},
}

func init() {
	sqlcmd.AddCommand(sql2structCmd)
	sql2structCmd.Flags().StringVarP(&username, "username", "u", "", "请输入数据库账号")
	sql2structCmd.Flags().StringVarP(&password, "password", "p", "", "请输入数据库密码")
	sql2structCmd.Flags().StringVarP(&host, "host", "", "127.0.0.1:3306", "请输入主机地址")
	sql2structCmd.Flags().StringVarP(&charset, "charset", "c", "utf8mb4", "请输入编码格式")
	sql2structCmd.Flags().StringVarP(&dbtype, "dbtype", "", "mysql", "请输入数据库类型")
	sql2structCmd.Flags().StringVarP(&dbname, "dbname", "", "", "请输入数据库名称")
	sql2structCmd.Flags().StringVarP(&tablename, "table", "", "", "请输入表格名称")
}
