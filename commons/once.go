/**
 * @Time    :2022/10/13 11:38
 * @Author  :Xiaoyu.Zhang
 */

package excommons

import (
	"sync"
	"sync/atomic"
)

// ErrOnce 出现err不更新锁的once
type ErrOnce struct {
	m    sync.Mutex
	done uint32
}

// Do
/**
*  @Description: 传入的函数f有返回值error，如果初始化失败，需要返回失败的error
*  @receiver o
*  @param f 需要单例执行的函数
*  @return error
 */
func (o *ErrOnce) Do(f func() error) error {
	if atomic.LoadUint32(&o.done) == 1 { //fast path
		return nil
	}
	return o.slowDo(f)
}

// 如果还没有初始化
func (o *ErrOnce) slowDo(f func() error) error {
	o.m.Lock()
	defer o.m.Unlock()
	var err error
	// 双检查，还没有初始化
	if o.done == 0 {
		err = f()
		if err == nil {
			// 初始化成功才将标记置为已初始化
			atomic.StoreUint32(&o.done, 1)
		}
	}
	return err
}
