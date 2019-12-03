package aoi

//Scene desc
type Scene struct {
	_head  *DEntry
	_tail  *DEntry
	_count int
}

//Initial desc
//@Method Initial desc:
func (slf *Scene) Initial() {
	slf._head = &DEntry{}
	slf._tail = &DEntry{}
	slf._head._xNext = slf._tail
	slf._head._yNext = slf._tail
	slf._tail._xPrev = slf._head
	slf._tail._yPrev = slf._head
}

//Enter desc
//@Method Enter desc: enter scene
//@Param (uint32) object id
//@Param (float64) x
//@Param (float64) y
//@Return (*DEntry) return node
func (slf *Scene) Enter(object uint32, x, y float64) *DEntry {
	n := &DEntry{_vKey: object}
	n._vPos.SetX(x)
	n._vPos.SetY(y)

	slf.add(n)

	return n
}

//Leave desc
//@Method Leave desc: leave object
func (slf *Scene) Leave(entry *DEntry) {
	if entry._xPrev == nil ||
		entry._xNext == nil ||
		entry._yPrev == nil ||
		entry._yNext == nil {
		return
	}

	//TODO: 找找周围的人，并发送事件

	entry._xPrev._xNext = entry._xNext
	entry._xNext._xPrev = entry._xPrev
	entry._yPrev._yNext = entry._yNext
	entry._yNext._yPrev = entry._yPrev

	entry._xPrev = nil
	entry._xNext = nil
	entry._yPrev = nil
	entry._yNext = nil
}

//Move desc
func (slf *Scene) Move(entry *DEntry, x, y float64) {

}

/*
 cout << "Cur is: " << node->key  << "（" << node ->x << "," << node ->y << ")" << endl;
        cout << "Print AOI:" << endl;

        // 往后找
        DoubleNode * cur = node->xNext;
        while(cur!=_tail)
        {
            if((cur->x - node->x) > xAreaLen)
            {
                break;
            }
            else
            {
                int inteval = 0;
                inteval = node->y - cur->y;
                if(inteval >= -yAreaLen && inteval <= yAreaLen)
                {
                    cout << "\t" << cur->key  << "(" << cur ->x << "," << cur ->y << ")" << endl;
                }
            }
            cur = cur->xNext;
        }

        // 往前找
        cur = node->xPrev;
        while(cur!=_head)
        {
            if((node->x - cur->x) > xAreaLen)
            {
                break;
            }
            else
            {
                int inteval = 0;
                inteval = node->y - cur->y;
                if(inteval >= -yAreaLen && inteval <= yAreaLen)
                {
                    cout << "\t" << cur->key  << "(" << cur ->x << "," << cur ->y << ")" << endl;
                }
            }
            cur = cur->xPrev;
        }*/

func (slf *Scene) add(n *DEntry) {
	//x
	cur := slf._head._xNext
	for cur != nil {
		if (cur._vPos.GetX() > n._vPos.GetX()) || (cur == slf._tail) {
			n._xNext = cur
			n._xPrev = cur._xPrev
			cur._xPrev._xNext = n
			cur._xPrev = n
			break
		}
		cur = cur._xNext
	}

	//y
	cur = slf._head._yNext
	for cur != nil {
		if (cur._vPos.GetY() > n._vPos.GetY()) || (cur == slf._tail) {
			n._yNext = cur
			n._yPrev = cur._yPrev
			cur._yPrev._yNext = n
			cur._yPrev = n
			break
		}
		cur = cur._yNext
	}

}
