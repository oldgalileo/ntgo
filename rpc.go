package ntgo

import "io"

type ValueRPC struct {
	DefVersion byte
	ProcedureName ValueString
	ParamSize uint8
	Params []RPCParam
	OutputSize uint8
	Outputs []RPCOutput
}

type RPCParam struct {
	Type EntryType
	Name ValueString
	DefaultVal EntryValue
}

type RPCOutput struct {
	Type EntryType
	Name ValueString
}

func DecodeRPC(r io.Reader) (*ValueRPC, error) {
	versionRaw := make([]byte, 1)
	_, versionErr := r.Read(versionRaw)
	if versionErr != nil {
		return nil, versionErr
	}
	procName, nameErr := DecodeString(r)
	if nameErr != nil {
		return nil, nameErr
	}
	paramSizeRaw := make([]byte, 1)
	_, paramSizeErr := r.Read(paramSizeRaw)
	if paramSizeErr != nil {
		return nil, paramSizeErr
	}
	paramSize := uint8(paramSizeRaw[0])
	params := make([]RPCParam, paramSize)
	for i := uint8(0); i < paramSize; i++ {
		param, paramErr := DecodeRPCParam(r)
		if paramErr != nil {
			return nil, paramErr
		}
		params[i] = param
	}
	outputSizeRaw := make([]byte, 1)
	_, outputSizeErr := r.Read(outputSizeRaw)
	if outputSizeErr != nil {
		return nil, outputSizeErr
	}
	outputSize := uint8(outputSizeRaw[0])
	outputs := make([]RPCOutput, outputSize)
	for i := uint8(0); i < outputSize; i++ {
		output, outputErr := DecodeRPCOutput(r)
		if outputErr != nil {
			return nil, outputErr
		}
		outputs[i] = output
	}
	return &ValueRPC{
		DefVersion: versionRaw[0],
		ProcedureName: *procName,
		ParamSize: paramSize,
		Params: params,
		OutputSize: outputSize,
		Outputs: outputs,
	}, nil
}

func (rpc *ValueRPC) GetRaw() []byte {
	raw := []byte{rpc.DefVersion}
	raw = append(raw, rpc.ProcedureName.GetRaw()...)
	raw = append(raw, byte(rpc.ParamSize))
	paramsRaw := []byte{}
	for i := uint8(0); i < rpc.ParamSize; i++ {
		paramsRaw = append(paramsRaw, rpc.Params[i].GetRaw()...)
	}
	raw = append(raw, paramsRaw...)
	raw = append(raw, byte(rpc.OutputSize))
	outputsRaw := []byte{}
	for i := uint8(0); i < rpc.OutputSize; i++ {
		outputsRaw = append(paramsRaw, rpc.Outputs[i].GetRaw()...)
	}
	raw = append(raw, outputsRaw...)
	return raw
}

func DecodeRPCParam(r io.Reader) (RPCParam, error) {
	rpcParam := RPCParam{}
	entryType, typeErr := DecodeEntryType(r)
	if typeErr != nil {
		return rpcParam, typeErr
	}
	rpcParam.Type = entryType
	name, nameErr := DecodeString(r)
	if nameErr != nil {
		return rpcParam, nameErr
	}
	rpcParam.Name = *name
	val, valErr := DecodeEntryWithType(r, entryType)
	if valErr != nil {
		return rpcParam, valErr
	}
	rpcParam.DefaultVal = val
	return rpcParam, nil
}

func (param RPCParam) GetRaw() []byte {
	raw := []byte{byte(param.Type)}
	raw = append(raw, param.Name.GetRaw()...)
	raw = append(raw, param.DefaultVal.GetRaw()...)
	return raw
}

func DecodeRPCOutput(r io.Reader) (RPCOutput, error) {
	output := RPCOutput{}
	entryType, typeErr := DecodeEntryType(r)
	if typeErr != nil {
		return output, typeErr
	}
	output.Type = entryType
	name, nameErr := DecodeString(r)
	if nameErr != nil {
		return output, nameErr
	}
	output.Name = *name
	return output, nil
}

func (output RPCOutput) GetRaw() []byte {
	raw := []byte{
		byte(output.Type),
	}
	raw = append(raw, output.Name.GetRaw()...)
	return raw
}