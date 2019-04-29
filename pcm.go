package main

func PCMF32TOI32(f32Buf []float32) (i32Buf []int32) {
	for _, v := range f32Buf {
		sample := v * 2147483647
		if sample > 2147483647 {
			sample = 2147483647
		}
		if sample < -2147483647 {
			sample = -2147483647
		}
		i32Buf = append(i32Buf, int32(sample))
	}
	return
}

func PCMI32TOInt(i32Buf []int32) (intBuf []int) {
	for _, v := range i32Buf {
		intBuf = append(intBuf, int(v))
	}
	return
}

func PCMF32ToI16(f32Buf []float32) (i16Buf []int16) {
	for _, v := range f32Buf {
		sample := v * 32767
		if sample > 32767 {
			sample = 32767
		}
		if sample < -32767 {
			sample = -32767
		}
		i16Buf = append(i16Buf, int16(sample))
	}
	return
}

func PCMI16TOInt(i16Buf []int16) (intBuf []int) {
	for _, v := range i16Buf {
		intBuf = append(intBuf, int(v))
	}
	return
}
