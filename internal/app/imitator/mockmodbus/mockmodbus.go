package mockmodbus

type QueryParams struct {
	Address, Quantity uint16
	Value             []byte
}

type MockModbus struct {
	WriteMultipleCoilsQueries     []QueryParams
	WriteMultipleRegistersQueries []QueryParams
}

func New() *MockModbus {
	return &MockModbus{}
}

func (m *MockModbus) ReadCoils(address, quantity uint16) (results []byte, err error) {
	return nil, nil
}

func (m *MockModbus) ReadDiscreteInputs(address, quantity uint16) (results []byte, err error) {
	return nil, nil
}

func (m *MockModbus) WriteSingleCoil(address, value uint16) (results []byte, err error) {
	return nil, nil
}

func (m *MockModbus) WriteMultipleCoils(address, quantity uint16, value []byte) (results []byte, err error) {
	m.WriteMultipleCoilsQueries = append(m.WriteMultipleCoilsQueries,
		QueryParams{
			Address:  address,
			Quantity: quantity,
			Value:    value,
		},
	)
	return nil, nil
}

func (m *MockModbus) ReadInputRegisters(address, quantity uint16) (results []byte, err error) {
	return nil, nil
}

func (m *MockModbus) ReadHoldingRegisters(address, quantity uint16) (results []byte, err error) {
	return nil, nil
}

func (m *MockModbus) WriteSingleRegister(address, value uint16) (results []byte, err error) {
	return nil, nil
}

func (m *MockModbus) WriteMultipleRegisters(address, quantity uint16, value []byte) (results []byte, err error) {
	m.WriteMultipleRegistersQueries = append(m.WriteMultipleRegistersQueries,
		QueryParams{
			Address:  address,
			Quantity: quantity,
			Value:    value,
		},
	)
	return nil, nil
}

func (m *MockModbus) ReadWriteMultipleRegisters(readAddress, readQuantity, writeAddress, writeQuantity uint16, value []byte) (results []byte, err error) {
	return nil, nil
}

func (m *MockModbus) MaskWriteRegister(address, andMask, orMask uint16) (results []byte, err error) {
	return nil, nil
}

func (m *MockModbus) ReadFIFOQueue(address uint16) (results []byte, err error) {
	return nil, nil
}

func (m *MockModbus) GetWriteMultipleCoilsQueries() []QueryParams {
	return m.WriteMultipleCoilsQueries
}

func (m *MockModbus) GetWriteMultipleRegistersQueries() []QueryParams {
	return m.WriteMultipleRegistersQueries
}
