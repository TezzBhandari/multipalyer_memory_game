"use client";

import React, { useState, useCallback, useEffect } from "react";
import useWebSocket, { ReadyState } from "react-use-websocket";

interface IJoinResponse {
  msgType: 1;
  data: {
    roomId: number;
    playerId: number;
  };
}

interface IGameStateResponse {
  msgType: 2;
  data: {
    counter: number;
    turn: number;
  };
}

interface IClientMessage {
  msgType: number;
  playerId: number;
  roomId: number;
}

export default function CounterGame() {
  //Public API that will echo messages sent to it back to the client

  const { sendJsonMessage, readyState, lastJsonMessage } = useWebSocket<
    IJoinResponse | IGameStateResponse
  >("ws://localhost:42069/game", { share: true });
  const [id, setId] = useState({
    roomId: -1,
    playerId: -1,
  });

  useEffect(() => {
    if (lastJsonMessage?.msgType === 1) {
      setId(lastJsonMessage.data);
    }
  }, [lastJsonMessage]);

  const handleClickSendMessage = useCallback(
    (type: "inc" | "dec") => {
      const payload: IClientMessage = {
        msgType: type === "inc" ? 1 : 2,
        roomId: id.roomId,
        playerId: id.playerId,
      };
      sendJsonMessage(payload);
    },
    [id],
  );

  const connectionStatus = {
    [ReadyState.CONNECTING]: "Connecting",
    [ReadyState.OPEN]: "Open",
    [ReadyState.CLOSING]: "Closing",
    [ReadyState.CLOSED]: "Closed",
    [ReadyState.UNINSTANTIATED]: "Uninstantiated",
  }[readyState];

  return (
    <>
      <div>
        <p>roomId: {id.roomId}</p>
        <p>playerId: {id.playerId}</p>
      </div>
      <div>
        <span>The WebSocket is currently {connectionStatus}</span>
        <span>Last message: {JSON.stringify(lastJsonMessage?.data)}</span>
      </div>
      <div className="flex gap-8">
        <button
          onClick={() => handleClickSendMessage("inc")}
          disabled={readyState !== ReadyState.OPEN}
          className="px-6 py-2 rounded-md capitalize bg-sky-400 text-2xl font-semibold text-white"
        >
          Increment
        </button>

        <button
          onClick={() => handleClickSendMessage("dec")}
          disabled={readyState !== ReadyState.OPEN}
          className="px-6 py-2 rounded-md capitalize bg-sky-400 text-2xl font-semibold text-white"
        >
          Decrement
        </button>
      </div>
    </>
  );
}
