package usr

import (
	"onward-path/internal/xui"
)

type User struct {
    xui.Client
}

func buyConfig(addClientParam *AddClientParam) {
    xuiAddClientParam := xui.AddClientRequestExternalAPI{
            Server: addClientParam.Server,
            ID: INBOUND_ID,
            Settings: xui.SettingsDecoded{
                Clients: []xui.ClientParam{
                    {
                        Email: addClientParam.Email,
                        Flow: addClientParam.Flow,
                        TotalGB: addClientParam.Total,
                        ExpiryTime: addClientParam.ExpiryTime,
                    },
                    },
            },
        }

   xui.AddClientInternal(xuiAddClientParam)
}
