# Features

Here are the features that are currently supported by yutu 🟢, and the ones that are planned to be supported in the future 🟡. The quota costs for each feature is also mentioned since there is a quota limits of 10,000 units/day.

<table>
  <thead>
    <tr>
      <th>Resource</th>
      <th>Action</th>
      <th>Quota</th>
      <th>MCP</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>abuseReports</td>
      <td>🟡 insert</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td>activities</td>
      <td>🟢 list</td>
      <td>1</td>
      <td>tool</td>
    </tr>
    <tr>
      <td rowspan="5">captions</td>
      <td>🟢 list</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 download</td>
      <td>?</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 insert</td>
      <td>400</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 update</td>
      <td>450</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 delete</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td rowspan="2">channels</td>
      <td>🟢 list</td>
      <td>1</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 update</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>channelBanners</td>
      <td>🟢 insert</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td rowspan="4">channelSections</td>
      <td>🟢 list</td>
      <td>1</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟡 insert</td>
      <td>50</td>
      <td></td>
    </tr>
    <tr>
      <td>🟡 update</td>
      <td>50</td>
      <td></td>
    </tr>
    <tr>
      <td>🟢 delete</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td rowspan="5">comments</td>
      <td>🟢 list</td>
      <td>1</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 insert</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 update</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 setModerationStatus</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 delete</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td rowspan="2">commentThreads</td>
      <td>🟢 list</td>
      <td>1</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 insert</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td><s>guideCategories</s></td>
      <td>🔴 <s>list deprecated API</s></td>
      <td>1</td>
      <td></td>
    </tr>
    <tr>
      <td>i18nLanguages</td>
      <td>🟢 list</td>
      <td>1</td>
      <td>resource</td>
    </tr>
    <tr>
      <td>i18nRegions</td>
      <td>🟢 list</td>
      <td>1</td>
      <td>resource</td>
    </tr>
    <tr>
      <td rowspan="7">liveBroadcasts</td>
      <td>🟡 list</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td>🟡 insert</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td>🟡 insertCuepoint</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td>🟡 update</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td>🟡 bind</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td>🟡 transition</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td>🟡 delete</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td rowspan="2">liveChatBans</td>
      <td>🟡 insert</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td>🟡 delete</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td rowspan="4">liveChatMessages</td>
      <td>🟡 list</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td>🟡 insert</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td>🟡 delete</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td>🟡 transition</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td rowspan="3">liveChatModerators</td>
      <td>🟡 list</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td>🟡 insert</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td>🟡 delete</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td rowspan="4">liveStreams</td>
      <td>🟡 list</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td>🟡 insert</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td>🟡 update</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td>🟡 delete</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td>members</td>
      <td>🟠 list <a href="https://github.com/eat-pray-ai/yutu/issues/3">🚫issue #3</a></td>
      <td>1</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>membershipsLevels</td>
      <td>🟠 list <a href="https://github.com/eat-pray-ai/yutu/issues/3">🚫issue #3</a></td>
      <td>1</td>
      <td>tool</td>
    </tr>
    <tr>
      <td rowspan="4">playlists</td>
      <td>🟢 list</td>
      <td>1</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 insert</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 update</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 delete</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td rowspan="4">playlistItems</td>
      <td>🟢 list</td>
      <td>1</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 insert</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 update</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 delete</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td rowspan="4">playlistImages</td>
      <td>🟢 list</td>
      <td>?</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 insert</td>
      <td>?</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 upload</td>
      <td>?</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 delete</td>
      <td>?</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>search</td>
      <td>🟢 list</td>
      <td>100</td>
      <td>tool</td>
    </tr>
    <tr>
      <td rowspan="3">subscriptions</td>
      <td>🟢 list</td>
      <td>1</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 insert</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 delete</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>superChatEvents</td>
      <td>🟢 list</td>
      <td>?</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>tests</td>
      <td>🟡 insert</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td rowspan="4">thirdPartyLinks</td>
      <td>🟡 list</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td>🟡 insert</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td>🟡 update</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td>🟡 delete</td>
      <td>?</td>
      <td></td>
    </tr>
    <tr>
      <td>thumbnails</td>
      <td>🟢 set</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td rowspan="7">videos</td>
      <td>🟢 list</td>
      <td>1</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 insert</td>
      <td>1600</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 update</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 rate</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 getRating</td>
      <td>1</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 reportAbuse</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 delete</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>videoAbuseReportReasons</td>
      <td>🟢 list</td>
      <td>1</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>videoCategories</td>
      <td>🟢 list</td>
      <td>1</td>
      <td>tool</td>
    </tr>
    <tr>
      <td rowspan="2">watermarks</td>
      <td>🟢 set</td>
      <td>50</td>
      <td>tool</td>
    </tr>
    <tr>
      <td>🟢 unset</td>
      <td>50</td>
      <td>tool</td>
    </tr>
  </tbody>
</table>
