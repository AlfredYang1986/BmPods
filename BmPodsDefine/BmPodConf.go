package BmPodsDefine

type Conf struct {
	Storages []struct{
		Name string
		Method string
		Daemons []string
	}
	Resources []struct{
		Name string
		Method string
		Storages []string
	}
	Models []string
	Services []struct{
		Name string
		Model string
		Resource string
		Storage []string
		Relationships struct{
			one2one []struct{
				Name string
				Method map[string]string
			}
			one2many []struct{
				Name string
				Method map[string]string
			}
		}
	}
	Daemons []struct{
		Name string
		Method string
		Args map[string]string
	}
}