const fs = require("fs");
require("dotenv").config();

const resourceMap = {
  0: "wheat",
  1: "fish",
  2: "iron",
  3: "wood",
  4: "coal",
  5: "oil",
  6: "gas",
};

const flattenedMarketList = (marketList) => {
  let res = [];
  marketList.forEach((listGroup) => {
    // first index holds metadata we don't need (e.g. "ultshared.UltOrder")
    res = res.concat(listGroup[1]);
  });
  return res;
};

const writePlayerFile = (payload) => {
  let players = payload.result.states[1].players;
  players = Object.values(players).map((p) => {
    return { ...p, isBot: p.lastLogin === 0 };
  });
  fs.writeFile("players.json", JSON.stringify(players, null, 4), () => {});
};

const writeMarketFile = (payload) => {
  const players = payload.result.states[1].players;
  const asks = payload.result.states[4].asks;
  const bids = payload.result.states[4].bids;
  const flatAsks = flattenedMarketList(asks);
  const flatBids = flattenedMarketList(bids);
  const marketData = [...flatAsks, ...flatBids].map((obj) => {
    return {
      ...obj,
      nation: players[obj.playerID]?.nationName ?? "",
      resource: resourceMap[obj.resourceType],
      isBot: players[obj.playerID]?.lastLogin === 0,
    };
  });
  fs.writeFile("market.json", JSON.stringify(marketData, null, 4), () => {});
};

const getSupremacyData = async () => {
  const res = await fetch("https://xgs-as-fwnq.c.bytro.com", {
    method: "POST",
    body: `{"@c":"ultshared.action.UltUpdateGameStateAction","version":"${process.env.VERSION}","client":"s1914-client-ultimate","adminLevel":0,"gameID":${process.env.GAME_ID},"playerID":0,"rights":"chat"}`,
    headers: {
      "Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
    },
  });

  const payload = await res.json();
  writePlayerFile(payload);
  writeMarketFile(payload);
};

getSupremacyData();
