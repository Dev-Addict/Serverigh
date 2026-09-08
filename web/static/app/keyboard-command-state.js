export const createCommandState = (setKeyStatus) => {
  let countBuffer = "";
  let commandBuffer = "";
  let timeout;

  const reset = (displayValue = "") => {
    countBuffer = "";
    commandBuffer = "";
    window.clearTimeout(timeout);
    setKeyStatus(displayValue);
  };

  const pending = () => {
    window.clearTimeout(timeout);
    setKeyStatus(`${countBuffer}${commandBuffer}`);
    timeout = window.setTimeout(() => {
      reset();
    }, 1200);
  };

  return {
    addCommand(key) {
      commandBuffer = `${commandBuffer}${key}`;
      pending();
    },
    addCount(key) {
      countBuffer = `${countBuffer}${key}`;
      pending();
    },
    count() {
      const count = Number.parseInt(countBuffer, 10);

      return Number.isNaN(count) || count < 1 ? 1 : count;
    },
    countPrefix() {
      return countBuffer;
    },
    hasCommand() {
      return commandBuffer.length > 0;
    },
    hasCount() {
      return countBuffer.length > 0;
    },
    reset,
    value() {
      return `${countBuffer}${commandBuffer}`;
    },
  };
};
