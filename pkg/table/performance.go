package table

import (
	"tradovatedataimport/pkg/csvdata"
	"tradovatedataimport/pkg/db"
	"tradovatedataimport/pkg/funcs"
)

func Performance() *Info {
	return NewInfo(
		"performance",
		Column{
			InputColumn: csvdata.NewColumn("symbol"),
			DbColumn:    db.NewPrimaryKeyColumn("symbol", db.TYPE_STRING),
		},
		Column{
			InputColumn: csvdata.NewColumn("_priceFormat"),
			DbColumn:    db.NewColumn("priceFormat", db.TYPE_INT),
		},
		Column{
			InputColumn: csvdata.NewColumn("_priceFormatType"),
			DbColumn:    db.NewColumn("priceFormatType", db.TYPE_INT),
		},
		Column{
			InputColumn: csvdata.NewColumn("_tickSize"),
			DbColumn:    db.NewColumn("tickSize", db.TYPE_DOUBLE),
		},
		Column{
			InputColumn: csvdata.NewColumn("buyFillId"),
			DbColumn:    db.NewPrimaryKeyColumn("buyFillId", db.TYPE_STRING),
		},
		Column{
			InputColumn: csvdata.NewColumn("sellFillId"),
			DbColumn:    db.NewPrimaryKeyColumn("sellFillId", db.TYPE_STRING),
		},
		Column{
			InputColumn: csvdata.NewColumn("qty"),
			DbColumn:    db.NewColumn("quantity", db.TYPE_DOUBLE),
		},
		Column{
			InputColumn: csvdata.NewColumn("buyPrice"),
			DbColumn:    db.NewColumn("buyPrice", db.TYPE_DOUBLE),
		},
		Column{
			InputColumn: csvdata.NewColumn("sellPrice"),
			DbColumn:    db.NewColumn("sellPrice", db.TYPE_DOUBLE),
		},
		Column{
			InputColumn: csvdata.NewColumn("pnl", funcs.RemoveNegativeParensFromCurrency),
			DbColumn:    db.NewColumn("pnl", db.TYPE_DOUBLE),
		},
		Column{
			InputColumn: csvdata.NewColumn("boughtTimestamp", funcs.CleanTimestamp),
			DbColumn:    db.NewColumn("boughtTimestamp", db.TYPE_DATETIME),
		},
		Column{
			InputColumn: csvdata.NewColumn("soldTimestamp", funcs.CleanTimestamp),
			DbColumn:    db.NewColumn("soldTimestamp", db.TYPE_DATETIME),
		},
		Column{
			InputColumn: csvdata.NewColumn("duration", funcs.TrimSpaces, funcs.CleanDurationAsSeconds),
			DbColumn:    db.NewColumn("durationSeconds", db.TYPE_UNSIGNED_BIG_INT),
		},
	)
}
