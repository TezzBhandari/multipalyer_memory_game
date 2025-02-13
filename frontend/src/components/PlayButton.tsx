"use client";
import { useRouter } from "next/navigation";
import React from "react";

export default function PlayButton() {
  const router = useRouter();

  const playHandler = () => {
    router.push("/game");
  };

  return (
    <button
      className="px-6 py-2 rounded-md capitalize bg-sky-400 text-2xl font-semibold text-white"
      onClick={playHandler}
    >
      play
    </button>
  );
}
