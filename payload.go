/**
 * @Author: fuxiao
 * @Author: 576101059@qq.com
 * @Date: 2021/7/19 16:24
 * @Desc: TODO
 */

package eventbus

import "github.com/dobyte/event-bus/internal/conv"

type Payload struct {
	Value []byte
}

func (p *Payload) Scan(any interface{}) error {
	return conv.Scan(p.Value, any)
}

func (p *Payload) String() string {
	return conv.String(p.Value)
}