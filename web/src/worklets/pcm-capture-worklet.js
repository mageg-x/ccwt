class PcmCaptureProcessor extends AudioWorkletProcessor {
  process(inputs, outputs) {
    const input = inputs[0]
    if (input && input[0] && input[0].length) {
      // 发送单声道 PCM Float32 数据给主线程
      this.port.postMessage(input[0].slice(0))
    }

    // 输出静音，避免啸叫
    const output = outputs[0]
    if (output && output[0]) {
      output[0].fill(0)
    }

    return true
  }
}

registerProcessor('pcm-capture-processor', PcmCaptureProcessor)

