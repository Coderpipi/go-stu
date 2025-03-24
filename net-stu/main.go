package main

type (
	TcpHeader struct {
		SrcPort      int16 // 原端口
		DstPort      int16 // 目标端口
		SeqNum       int32 // 序列号
		AckNum       int32 // 确认号
		TcpHeaderLen int8  // TCP头部总长度 4位
		Resv         int8  // 4位
		Flags        struct {
			NS  int8 // 标志位
			CWR int8 // 标志位
			ECE int8 // 标志位
			URG int8 // 标志位
			ACK int8 // 标志位
			PSH int8 // 标志位
			RST int8 // 标志位
			SYN int8 // 标志位
		} // 标志位
		Window   int16 // 窗口大小
		Checksum int16 // 校验和
		Urgent   int16 // 紧急指针
	}

	IpHeader struct {
		Version    int8  // 4位
		HeaderLen  int8  // 头部长度 4位
		Tos        int8  // 服务类型
		TotalLen   int16 // 总长度
		Id         int16 // 标识
		Flags      int8  // 标志 3位
		FragOffset int8  // 片偏移 13位
		Ttl        int8  // 生存时间 8位
		Protocol   int8  // 协议 8位
		Checksum   int16 // 校验和 16位
		SrcAddr    int32 // 源地址
		DstAddr    int32 // 目标地址
	}

	UdpHeader struct {
		SrcPort  int16 // 源端口
		DstPort  int16 // 目标端口
		Length   int16 // 长度
		Checksum int16 // 校验和
	}
)

func main() {
	/*
				tcp 三次握手
				1. seq_num = x, syn = 1
			    2. ack = 1,syn =1, seq_num = y, ack_num = x + 1
		        3. ack = 1, ack_num = y + 1, seq_num = z
	*/
}
