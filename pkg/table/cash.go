package table

import (
	"tradovatedataimport/pkg/csvdata"
	"tradovatedataimport/pkg/db"
	"tradovatedataimport/pkg/funcs"
)

func Cash() *Info {
	return NewInfo(
		"cash",
		Column{
			InputColumn: csvdata.NewColumn("Account"),
			DbColumn:    db.NewPrimaryKeyColumn("account", db.TYPE_STRING),
		},
		Column{
			InputColumn: csvdata.NewColumn("Transaction ID"),
			DbColumn:    db.NewPrimaryKeyColumn("transactionId", db.TYPE_STRING),
		},
		Column{
			InputColumn: csvdata.NewColumn("Timestamp", funcs.CleanTimestamp),
			DbColumn:    db.NewColumn("timestamp", db.TYPE_DATETIME),
		},
		Column{
			InputColumn: csvdata.NewColumn("Date"),
			DbColumn:    db.NewColumn("date", db.TYPE_DATE),
		},
		Column{
			InputColumn: csvdata.NewColumn("Delta", funcs.RemoveCommas),
			DbColumn:    db.NewColumn("delta", db.TYPE_DOUBLE),
		},
		Column{
			InputColumn: csvdata.NewColumn("Amount", funcs.RemoveCommas),
			DbColumn:    db.NewColumn("amount", db.TYPE_DOUBLE),
		},
		Column{
			InputColumn: csvdata.NewColumn("Cash Change Type"),
			DbColumn:    db.NewColumn("cashChangeType", db.TYPE_STRING),
		},
		Column{
			InputColumn: csvdata.NewColumn("Currency"),
			DbColumn:    db.NewColumn("currency", db.TYPE_STRING),
		},
		Column{
			InputColumn: csvdata.NewColumn("Contract"),
			DbColumn:    db.NewColumn("contract", db.TYPE_STRING),
		},
	)
}
