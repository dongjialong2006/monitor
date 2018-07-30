# monitor
监控系统CPU和内存的使用情况

## 描述
实时打印当前进程占用CPU和内存的比率(用户可自行设置)

## 引用

- import "github.com/dongjialong2006/monitor"


## 使用
{
	
	import "time"
	import "github.com/dongjialong2006/monitor"
	
	
	func main() {
		monitor.Watch(context.Background())
		time.Sleep(5*tiem.Second)
		
		monitor.Stop()
	}
}
