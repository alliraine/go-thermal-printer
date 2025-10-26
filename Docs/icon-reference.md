# Material Icon Catalog for Templates

The template engine ships with the complete Material Design icon set from `golang.org/x/exp/shiny/materialdesign/icons`. Use the canonical names below with the `icon` helper to render scalable vector icons on the printer.


A ready-to-print demonstration lives in `templates/icons-demo.tmpl`. It prints a few icons with different sizing options.

## Using the `icon` helper

```gotemplate
{{icon "ActionAlarm"}}
{{icon "ActionAlarm" 128}}
{{icon "ActionAlarm" 128 24}}
```

* **Name** – Accepts any case-insensitive, alphanumeric, hyphenated, or underscored variant. `ActionAlarm`, `action-alarm`, and `ACTION_ALARM` all resolve to the same icon.
* **Width (optional)** – First numeric argument, in printer dots. Defaults to 96 dots (≈ 12 mm on a 203 dpi printer). Icons keep their aspect ratio, so height is computed automatically.
* **Feed (optional)** – Second numeric argument that feeds additional blank lines immediately after the icon. Values are clamped to 0–255 lines (default `0`).
* The helper resets the printer’s code page after each icon, so subsequent text prints normally.

---

## Icon Catalog

Below is the complete list of canonical icon names organised by category. Each entry can be copied directly into a template, for example `{{icon "MapsDirectionsBike" 120 8}}`.

## AV Icons
| Name | Name | Name | Name |
|------|------|------|------|
| AVAVTimer | AVAddToQueue | AVAirplay | AVAlbum |
| AVArtTrack | AVBrandingWatermark | AVCallToAction | AVClosedCaption |
| AVEqualizer | AVExplicit | AVFastForward | AVFastRewind |
| AVFeaturedPlayList | AVFeaturedVideo | AVFiberDVR | AVFiberManualRecord |
| AVFiberNew | AVFiberPin | AVFiberSmartRecord | AVForward10 |
| AVForward30 | AVForward5 | AVGames | AVHD |
| AVHearing | AVHighQuality | AVLibraryAdd | AVLibraryBooks |
| AVLibraryMusic | AVLoop | AVMic | AVMicNone |
| AVMicOff | AVMovie | AVMusicVideo | AVNewReleases |
| AVNotInterested | AVNote | AVPause | AVPauseCircleFilled |
| AVPauseCircleOutline | AVPlayArrow | AVPlayCircleFilled | AVPlayCircleOutline |
| AVPlaylistAdd | AVPlaylistAddCheck | AVPlaylistPlay | AVQueue |
| AVQueueMusic | AVQueuePlayNext | AVRadio | AVRecentActors |
| AVRemoveFromQueue | AVRepeat | AVRepeatOne | AVReplay |
| AVReplay10 | AVReplay30 | AVReplay5 | AVShuffle |
| AVSkipNext | AVSkipPrevious | AVSlowMotionVideo | AVSnooze |
| AVSortByAlpha | AVStop | AVSubscriptions | AVSubtitles |
| AVSurroundSound | AVVideoCall | AVVideoLabel | AVVideoLibrary |
| AVVideocam | AVVideocamOff | AVVolumeDown | AVVolumeMute |
| AVVolumeOff | AVVolumeUp | AVWeb | AVWebAsset |

## Action Icons
| Name | Name | Name | Name |
|------|------|------|------|
| Action3DRotation | ActionAccessibility | ActionAccessible | ActionAccountBalance |
| ActionAccountBalanceWallet | ActionAccountBox | ActionAccountCircle | ActionAddShoppingCart |
| ActionAlarm | ActionAlarmAdd | ActionAlarmOff | ActionAlarmOn |
| ActionAllOut | ActionAndroid | ActionAnnouncement | ActionAspectRatio |
| ActionAssessment | ActionAssignment | ActionAssignmentInd | ActionAssignmentLate |
| ActionAssignmentReturn | ActionAssignmentReturned | ActionAssignmentTurnedIn | ActionAutorenew |
| ActionBackup | ActionBook | ActionBookmark | ActionBookmarkBorder |
| ActionBugReport | ActionBuild | ActionCached | ActionCameraEnhance |
| ActionCardGiftcard | ActionCardMembership | ActionCardTravel | ActionChangeHistory |
| ActionCheckCircle | ActionChromeReaderMode | ActionClass | ActionCode |
| ActionCompareArrows | ActionCopyright | ActionCreditCard | ActionDNS |
| ActionDashboard | ActionDateRange | ActionDelete | ActionDeleteForever |
| ActionDescription | ActionDone | ActionDoneAll | ActionDonutLarge |
| ActionDonutSmall | ActionEject | ActionEuroSymbol | ActionEvent |
| ActionEventSeat | ActionExitToApp | ActionExplore | ActionExtension |
| ActionFace | ActionFavorite | ActionFavoriteBorder | ActionFeedback |
| ActionFindInPage | ActionFindReplace | ActionFingerprint | ActionFlightLand |
| ActionFlightTakeoff | ActionFlipToBack | ActionFlipToFront | ActionGIF |
| ActionGTranslate | ActionGavel | ActionGetApp | ActionGrade |
| ActionGroupWork | ActionHTTP | ActionHTTPS | ActionHelp |
| ActionHelpOutline | ActionHighlightOff | ActionHistory | ActionHome |
| ActionHourglassEmpty | ActionHourglassFull | ActionImportantDevices | ActionInfo |
| ActionInfoOutline | ActionInput | ActionInvertColors | ActionLabel |
| ActionLabelOutline | ActionLanguage | ActionLaunch | ActionLightbulbOutline |
| ActionLineStyle | ActionLineWeight | ActionList | ActionLock |
| ActionLockOpen | ActionLockOutline | ActionLoyalty | ActionMarkUnreadMailbox |
| ActionMotorcycle | ActionNoteAdd | ActionOfflinePin | ActionOpacity |
| ActionOpenInBrowser | ActionOpenInNew | ActionOpenWith | ActionPageview |
| ActionPanTool | ActionPayment | ActionPermCameraMic | ActionPermContactCalendar |
| ActionPermDataSetting | ActionPermDeviceInformation | ActionPermIdentity | ActionPermMedia |
| ActionPermPhoneMsg | ActionPermScanWiFi | ActionPets | ActionPictureInPicture |
| ActionPictureInPictureAlt | ActionPlayForWork | ActionPolymer | ActionPowerSettingsNew |
| ActionPregnantWoman | ActionPrint | ActionQueryBuilder | ActionQuestionAnswer |
| ActionReceipt | ActionRecordVoiceOver | ActionRedeem | ActionRemoveShoppingCart |
| ActionReorder | ActionReportProblem | ActionRestore | ActionRestorePage |
| ActionRoom | ActionRoundedCorner | ActionRowing | ActionSchedule |
| ActionSearch | ActionSettings | ActionSettingsApplications | ActionSettingsBackupRestore |
| ActionSettingsBluetooth | ActionSettingsBrightness | ActionSettingsCell | ActionSettingsEthernet |
| ActionSettingsInputAntenna | ActionSettingsInputComponent | ActionSettingsInputComposite | ActionSettingsInputHDMI |
| ActionSettingsInputSVideo | ActionSettingsOverscan | ActionSettingsPhone | ActionSettingsPower |
| ActionSettingsRemote | ActionSettingsVoice | ActionShop | ActionShopTwo |
| ActionShoppingBasket | ActionShoppingCart | ActionSpeakerNotes | ActionSpeakerNotesOff |
| ActionSpellcheck | ActionStarRate | ActionStars | ActionStore |
| ActionSubject | ActionSupervisorAccount | ActionSwapHoriz | ActionSwapVert |
| ActionSwapVerticalCircle | ActionSystemUpdateAlt | ActionTOC | ActionTab |
| ActionTabUnselected | ActionTheaters | ActionThumbDown | ActionThumbUp |
| ActionThumbsUpDown | ActionTimeline | ActionToday | ActionToll |
| ActionTouchApp | ActionTrackChanges | ActionTranslate | ActionTrendingDown |
| ActionTrendingFlat | ActionTrendingUp | ActionTurnedIn | ActionTurnedInNot |
| ActionUpdate | ActionVerifiedUser | ActionViewAgenda | ActionViewArray |
| ActionViewCarousel | ActionViewColumn | ActionViewDay | ActionViewHeadline |
| ActionViewList | ActionViewModule | ActionViewQuilt | ActionViewStream |
| ActionViewWeek | ActionVisibility | ActionVisibilityOff | ActionWatchLater |
| ActionWork | ActionYoutubeSearchedFor | ActionZoomIn | ActionZoomOut |

## Alert Icons
| Name | Name | Name | Name |
|------|------|------|------|
| AlertAddAlert | AlertError | AlertErrorOutline | AlertWarning |

## Communication Icons
| Name | Name | Name | Name |
|------|------|------|------|
| CommunicationBusiness | CommunicationCall | CommunicationCallEnd | CommunicationCallMade |
| CommunicationCallMerge | CommunicationCallMissed | CommunicationCallMissedOutgoing | CommunicationCallReceived |
| CommunicationCallSplit | CommunicationChat | CommunicationChatBubble | CommunicationChatBubbleOutline |
| CommunicationClearAll | CommunicationComment | CommunicationContactMail | CommunicationContactPhone |
| CommunicationContacts | CommunicationDialerSIP | CommunicationDialpad | CommunicationEmail |
| CommunicationForum | CommunicationImportContacts | CommunicationImportExport | CommunicationInvertColorsOff |
| CommunicationLiveHelp | CommunicationLocationOff | CommunicationLocationOn | CommunicationMailOutline |
| CommunicationMessage | CommunicationNoSIM | CommunicationPhone | CommunicationPhoneLinkErase |
| CommunicationPhoneLinkLock | CommunicationPhoneLinkRing | CommunicationPhoneLinkSetup | CommunicationPortableWiFiOff |
| CommunicationPresentToAll | CommunicationRSSFeed | CommunicationRingVolume | CommunicationScreenShare |
| CommunicationSpeakerPhone | CommunicationStayCurrentLandscape | CommunicationStayCurrentPortrait | CommunicationStayPrimaryLandscape |
| CommunicationStayPrimaryPortrait | CommunicationStopScreenShare | CommunicationSwapCalls | CommunicationTextSMS |
| CommunicationVPNKey | CommunicationVoicemail |  |  |

## Content Icons
| Name | Name | Name | Name |
|------|------|------|------|
| ContentAdd | ContentAddBox | ContentAddCircle | ContentAddCircleOutline |
| ContentArchive | ContentBackspace | ContentBlock | ContentClear |
| ContentContentCopy | ContentContentCut | ContentContentPaste | ContentCreate |
| ContentDeleteSweep | ContentDrafts | ContentFilterList | ContentFlag |
| ContentFontDownload | ContentForward | ContentGesture | ContentInbox |
| ContentLink | ContentLowPriority | ContentMail | ContentMarkUnread |
| ContentMoveToInbox | ContentNextWeek | ContentRedo | ContentRemove |
| ContentRemoveCircle | ContentRemoveCircleOutline | ContentReply | ContentReplyAll |
| ContentReport | ContentSave | ContentSelectAll | ContentSend |
| ContentSort | ContentTextFormat | ContentUnarchive | ContentUndo |
| ContentWeekend |  |  |  |

## Device Icons
| Name | Name | Name | Name |
|------|------|------|------|
| DeviceAccessAlarm | DeviceAccessAlarms | DeviceAccessTime | DeviceAddAlarm |
| DeviceAirplaneModeActive | DeviceAirplaneModeInactive | DeviceBattery20 | DeviceBattery30 |
| DeviceBattery50 | DeviceBattery60 | DeviceBattery80 | DeviceBattery90 |
| DeviceBatteryAlert | DeviceBatteryCharging20 | DeviceBatteryCharging30 | DeviceBatteryCharging50 |
| DeviceBatteryCharging60 | DeviceBatteryCharging80 | DeviceBatteryCharging90 | DeviceBatteryChargingFull |
| DeviceBatteryFull | DeviceBatteryStd | DeviceBatteryUnknown | DeviceBluetooth |
| DeviceBluetoothConnected | DeviceBluetoothDisabled | DeviceBluetoothSearching | DeviceBrightnessAuto |
| DeviceBrightnessHigh | DeviceBrightnessLow | DeviceBrightnessMedium | DeviceDVR |
| DeviceDataUsage | DeviceDeveloperMode | DeviceDevices | DeviceGPSFixed |
| DeviceGPSNotFixed | DeviceGPSOff | DeviceGraphicEq | DeviceLocationDisabled |
| DeviceLocationSearching | DeviceNFC | DeviceNetworkCell | DeviceNetworkWiFi |
| DeviceSDStorage | DeviceScreenLockLandscape | DeviceScreenLockPortrait | DeviceScreenLockRotation |
| DeviceScreenRotation | DeviceSettingsSystemDaydream | DeviceSignalCellular0Bar | DeviceSignalCellular1Bar |
| DeviceSignalCellular2Bar | DeviceSignalCellular3Bar | DeviceSignalCellular4Bar | DeviceSignalCellularConnectedNoInternet0Bar |
| DeviceSignalCellularConnectedNoInternet1Bar | DeviceSignalCellularConnectedNoInternet2Bar | DeviceSignalCellularConnectedNoInternet3Bar | DeviceSignalCellularConnectedNoInternet4Bar |
| DeviceSignalCellularNoSIM | DeviceSignalCellularNull | DeviceSignalCellularOff | DeviceSignalWiFi0Bar |
| DeviceSignalWiFi1Bar | DeviceSignalWiFi1BarLock | DeviceSignalWiFi2Bar | DeviceSignalWiFi2BarLock |
| DeviceSignalWiFi3Bar | DeviceSignalWiFi3BarLock | DeviceSignalWiFi4Bar | DeviceSignalWiFi4BarLock |
| DeviceSignalWiFiOff | DeviceStorage | DeviceUSB | DeviceWallpaper |
| DeviceWiFiLock | DeviceWiFiTethering | DeviceWidgets |  |

## Editor Icons
| Name | Name | Name | Name |
|------|------|------|------|
| EditorAttachFile | EditorAttachMoney | EditorBorderAll | EditorBorderBottom |
| EditorBorderClear | EditorBorderColor | EditorBorderHorizontal | EditorBorderInner |
| EditorBorderLeft | EditorBorderOuter | EditorBorderRight | EditorBorderStyle |
| EditorBorderTop | EditorBorderVertical | EditorBubbleChart | EditorDragHandle |
| EditorFormatAlignCenter | EditorFormatAlignJustify | EditorFormatAlignLeft | EditorFormatAlignRight |
| EditorFormatBold | EditorFormatClear | EditorFormatColorFill | EditorFormatColorReset |
| EditorFormatColorText | EditorFormatIndentDecrease | EditorFormatIndentIncrease | EditorFormatItalic |
| EditorFormatLineSpacing | EditorFormatListBulleted | EditorFormatListNumbered | EditorFormatPaint |
| EditorFormatQuote | EditorFormatShapes | EditorFormatSize | EditorFormatStrikethrough |
| EditorFormatTextDirectionLToR | EditorFormatTextDirectionRToL | EditorFormatUnderlined | EditorFunctions |
| EditorHighlight | EditorInsertChart | EditorInsertComment | EditorInsertDriveFile |
| EditorInsertEmoticon | EditorInsertInvitation | EditorInsertLink | EditorInsertPhoto |
| EditorLinearScale | EditorMergeType | EditorModeComment | EditorModeEdit |
| EditorMonetizationOn | EditorMoneyOff | EditorMultilineChart | EditorPieChart |
| EditorPieChartOutlined | EditorPublish | EditorShortText | EditorShowChart |
| EditorSpaceBar | EditorStrikethroughS | EditorTextFields | EditorTitle |
| EditorVerticalAlignBottom | EditorVerticalAlignCenter | EditorVerticalAlignTop | EditorWrapText |

## File Icons
| Name | Name | Name | Name |
|------|------|------|------|
| FileAttachment | FileCloud | FileCloudCircle | FileCloudDone |
| FileCloudDownload | FileCloudOff | FileCloudQueue | FileCloudUpload |
| FileCreateNewFolder | FileFileDownload | FileFileUpload | FileFolder |
| FileFolderOpen | FileFolderShared |  |  |

## Hardware Icons
| Name | Name | Name | Name |
|------|------|------|------|
| HardwareCast | HardwareCastConnected | HardwareComputer | HardwareDesktopMac |
| HardwareDesktopWindows | HardwareDeveloperBoard | HardwareDeviceHub | HardwareDevicesOther |
| HardwareDock | HardwareGamepad | HardwareHeadset | HardwareHeadsetMic |
| HardwareKeyboard | HardwareKeyboardArrowDown | HardwareKeyboardArrowLeft | HardwareKeyboardArrowRight |
| HardwareKeyboardArrowUp | HardwareKeyboardBackspace | HardwareKeyboardCapslock | HardwareKeyboardHide |
| HardwareKeyboardReturn | HardwareKeyboardTab | HardwareKeyboardVoice | HardwareLaptop |
| HardwareLaptopChromebook | HardwareLaptopMac | HardwareLaptopWindows | HardwareMemory |
| HardwareMouse | HardwarePhoneAndroid | HardwarePhoneIPhone | HardwarePhoneLink |
| HardwarePhoneLinkOff | HardwarePowerInput | HardwareRouter | HardwareSIMCard |
| HardwareScanner | HardwareSecurity | HardwareSmartphone | HardwareSpeaker |
| HardwareSpeakerGroup | HardwareTV | HardwareTablet | HardwareTabletAndroid |
| HardwareTabletMac | HardwareToys | HardwareVideogameAsset | HardwareWatch |

## Image Icons
| Name | Name | Name | Name |
|------|------|------|------|
| ImageAddAPhoto | ImageAddToPhotos | ImageAdjust | ImageAssistant |
| ImageAssistantPhoto | ImageAudiotrack | ImageBlurCircular | ImageBlurLinear |
| ImageBlurOff | ImageBlurOn | ImageBrightness1 | ImageBrightness2 |
| ImageBrightness3 | ImageBrightness4 | ImageBrightness5 | ImageBrightness6 |
| ImageBrightness7 | ImageBrokenImage | ImageBrush | ImageBurstMode |
| ImageCamera | ImageCameraAlt | ImageCameraFront | ImageCameraRear |
| ImageCameraRoll | ImageCenterFocusStrong | ImageCenterFocusWeak | ImageCollections |
| ImageCollectionsBookmark | ImageColorLens | ImageColorize | ImageCompare |
| ImageControlPoint | ImageControlPointDuplicate | ImageCrop | ImageCrop169 |
| ImageCrop32 | ImageCrop54 | ImageCrop75 | ImageCropDIN |
| ImageCropFree | ImageCropLandscape | ImageCropOriginal | ImageCropPortrait |
| ImageCropRotate | ImageCropSquare | ImageDehaze | ImageDetails |
| ImageEdit | ImageExposure | ImageExposureNeg1 | ImageExposureNeg2 |
| ImageExposurePlus1 | ImageExposurePlus2 | ImageExposureZero | ImageFilter |
| ImageFilter1 | ImageFilter2 | ImageFilter3 | ImageFilter4 |
| ImageFilter5 | ImageFilter6 | ImageFilter7 | ImageFilter8 |
| ImageFilter9 | ImageFilter9Plus | ImageFilterBAndW | ImageFilterCenterFocus |
| ImageFilterDrama | ImageFilterFrames | ImageFilterHDR | ImageFilterNone |
| ImageFilterTiltShift | ImageFilterVintage | ImageFlare | ImageFlashAuto |
| ImageFlashOff | ImageFlashOn | ImageFlip | ImageGradient |
| ImageGrain | ImageGridOff | ImageGridOn | ImageHDROff |
| ImageHDROn | ImageHDRStrong | ImageHDRWeak | ImageHealing |
| ImageISO | ImageImage | ImageImageAspectRatio | ImageLandscape |
| ImageLeakAdd | ImageLeakRemove | ImageLens | ImageLinkedCamera |
| ImageLooks | ImageLooks3 | ImageLooks4 | ImageLooks5 |
| ImageLooks6 | ImageLooksOne | ImageLooksTwo | ImageLoupe |
| ImageMonochromePhotos | ImageMovieCreation | ImageMovieFilter | ImageMusicNote |
| ImageNature | ImageNaturePeople | ImageNavigateBefore | ImageNavigateNext |
| ImagePalette | ImagePanorama | ImagePanoramaFishEye | ImagePanoramaHorizontal |
| ImagePanoramaVertical | ImagePanoramaWideAngle | ImagePhoto | ImagePhotoAlbum |
| ImagePhotoCamera | ImagePhotoFilter | ImagePhotoLibrary | ImagePhotoSizeSelectActual |
| ImagePhotoSizeSelectLarge | ImagePhotoSizeSelectSmall | ImagePictureAsPDF | ImagePortrait |
| ImageRemoveRedEye | ImageRotate90DegreesCCW | ImageRotateLeft | ImageRotateRight |
| ImageSlideshow | ImageStraighten | ImageStyle | ImageSwitchCamera |
| ImageSwitchVideo | ImageTagFaces | ImageTexture | ImageTimeLapse |
| ImageTimer | ImageTimer10 | ImageTimer3 | ImageTimerOff |
| ImageTonality | ImageTransform | ImageTune | ImageViewComfy |
| ImageViewCompact | ImageVignette | ImageWBAuto | ImageWBCloudy |
| ImageWBIncandescent | ImageWBIridescent | ImageWBSunny |  |

## Maps Icons
| Name | Name | Name | Name |
|------|------|------|------|
| MapsAddLocation | MapsBeenhere | MapsDirections | MapsDirectionsBike |
| MapsDirectionsBoat | MapsDirectionsBus | MapsDirectionsCar | MapsDirectionsRailway |
| MapsDirectionsRun | MapsDirectionsSubway | MapsDirectionsTransit | MapsDirectionsWalk |
| MapsEVStation | MapsEditLocation | MapsFlight | MapsHotel |
| MapsLayers | MapsLayersClear | MapsLocalATM | MapsLocalActivity |
| MapsLocalAirport | MapsLocalBar | MapsLocalCafe | MapsLocalCarWash |
| MapsLocalConvenienceStore | MapsLocalDining | MapsLocalDrink | MapsLocalFlorist |
| MapsLocalGasStation | MapsLocalGroceryStore | MapsLocalHospital | MapsLocalHotel |
| MapsLocalLaundryService | MapsLocalLibrary | MapsLocalMall | MapsLocalMovies |
| MapsLocalOffer | MapsLocalParking | MapsLocalPharmacy | MapsLocalPhone |
| MapsLocalPizza | MapsLocalPlay | MapsLocalPostOffice | MapsLocalPrintshop |
| MapsLocalSee | MapsLocalShipping | MapsLocalTaxi | MapsMap |
| MapsMyLocation | MapsNavigation | MapsNearMe | MapsPersonPin |
| MapsPersonPinCircle | MapsPinDrop | MapsPlace | MapsRateReview |
| MapsRestaurant | MapsRestaurantMenu | MapsSatellite | MapsStoreMallDirectory |
| MapsStreetView | MapsSubway | MapsTerrain | MapsTraffic |
| MapsTrain | MapsTram | MapsTransferWithinAStation | MapsZoomOutMap |

## Navigation Icons
| Name | Name | Name | Name |
|------|------|------|------|
| NavigationApps | NavigationArrowBack | NavigationArrowDownward | NavigationArrowDropDown |
| NavigationArrowDropDownCircle | NavigationArrowDropUp | NavigationArrowForward | NavigationArrowUpward |
| NavigationCancel | NavigationCheck | NavigationChevronLeft | NavigationChevronRight |
| NavigationClose | NavigationExpandLess | NavigationExpandMore | NavigationFirstPage |
| NavigationFullscreen | NavigationFullscreenExit | NavigationLastPage | NavigationMenu |
| NavigationMoreHoriz | NavigationMoreVert | NavigationRefresh | NavigationSubdirectoryArrowLeft |
| NavigationSubdirectoryArrowRight | NavigationUnfoldLess | NavigationUnfoldMore |  |

## Notification Icons
| Name | Name | Name | Name |
|------|------|------|------|
| NotificationADB | NotificationAirlineSeatFlat | NotificationAirlineSeatFlatAngled | NotificationAirlineSeatIndividualSuite |
| NotificationAirlineSeatLegroomExtra | NotificationAirlineSeatLegroomNormal | NotificationAirlineSeatLegroomReduced | NotificationAirlineSeatReclineExtra |
| NotificationAirlineSeatReclineNormal | NotificationBluetoothAudio | NotificationConfirmationNumber | NotificationDiscFull |
| NotificationDoNotDisturb | NotificationDoNotDisturbAlt | NotificationDoNotDisturbOff | NotificationDoNotDisturbOn |
| NotificationDriveETA | NotificationEnhancedEncryption | NotificationEventAvailable | NotificationEventBusy |
| NotificationEventNote | NotificationFolderSpecial | NotificationLiveTV | NotificationMMS |
| NotificationMore | NotificationNetworkCheck | NotificationNetworkLocked | NotificationNoEncryption |
| NotificationOnDemandVideo | NotificationPersonalVideo | NotificationPhoneBluetoothSpeaker | NotificationPhoneForwarded |
| NotificationPhoneInTalk | NotificationPhoneLocked | NotificationPhoneMissed | NotificationPhonePaused |
| NotificationPower | NotificationPriorityHigh | NotificationRVHookup | NotificationSDCard |
| NotificationSIMCardAlert | NotificationSMS | NotificationSMSFailed | NotificationSync |
| NotificationSyncDisabled | NotificationSyncProblem | NotificationSystemUpdate | NotificationTapAndPlay |
| NotificationTimeToLeave | NotificationVPNLock | NotificationVibration | NotificationVoiceChat |
| NotificationWC | NotificationWiFi |  |  |

## Places Icons
| Name | Name | Name | Name |
|------|------|------|------|
| PlacesACUnit | PlacesAirportShuttle | PlacesAllInclusive | PlacesBeachAccess |
| PlacesBusinessCenter | PlacesCasino | PlacesChildCare | PlacesChildFriendly |
| PlacesFitnessCenter | PlacesFreeBreakfast | PlacesGolfCourse | PlacesHotTub |
| PlacesKitchen | PlacesPool | PlacesRVHookup | PlacesRoomService |
| PlacesSmokeFree | PlacesSmokingRooms | PlacesSpa |  |

## Social Icons
| Name | Name | Name | Name |
|------|------|------|------|
| SocialCake | SocialDomain | SocialGroup | SocialGroupAdd |
| SocialLocationCity | SocialMood | SocialMoodBad | SocialNotifications |
| SocialNotificationsActive | SocialNotificationsNone | SocialNotificationsOff | SocialNotificationsPaused |
| SocialPages | SocialPartyMode | SocialPeople | SocialPeopleOutline |
| SocialPerson | SocialPersonAdd | SocialPersonOutline | SocialPlusOne |
| SocialPoll | SocialPublic | SocialSchool | SocialSentimentDissatisfied |
| SocialSentimentNeutral | SocialSentimentSatisfied | SocialSentimentVeryDissatisfied | SocialSentimentVerySatisfied |
| SocialShare | SocialWhatsHot |  |  |

## Toggle Icons
| Name | Name | Name | Name |
|------|------|------|------|
| ToggleCheckBox | ToggleCheckBoxOutlineBlank | ToggleIndeterminateCheckBox | ToggleRadioButtonChecked |
| ToggleRadioButtonUnchecked | ToggleStar | ToggleStarBorder | ToggleStarHalf |
