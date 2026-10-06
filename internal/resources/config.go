// Copyright (C) 2026 Paul Dwerryhouse <paul@dwerryhouse.com.au>
//
// This file is part of terraform-provider-immich.
//
// terraform-provider-immich is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// terraform-provider-immich is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with terraform-provider-immich.  If not, see <https://www.gnu.org/licenses/>.

package resources

import (
	"context"
	"fmt"

	immichclient "codeberg.org/pdwerryhouse/immich-client-go/client"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &configResource{}
	_ resource.ResourceWithConfigure   = &configResource{}
	_ resource.ResourceWithImportState = &configResource{}
)

func NewConfigResource() resource.Resource {
	return &configResource{}
}

type configResource struct {
	client *immichclient.Client
}

type databaseBackupConfigModel struct {
	CronExpression types.String `tfsdk:"cron_expression"`
	Enabled        types.Bool   `tfsdk:"enabled"`
	KeepLastAmount types.Int64  `tfsdk:"keep_last_amount"`
}

type BackupsConfigModel struct {
	Database databaseBackupConfigModel `tfsdk:"database"`
}

type ffmpegRealtimeConfigModel struct {
	Enabled     types.Bool `tfsdk:"enabled"`
	Resolutions types.List `tfsdk:"resolutions"`
	VideoCodecs types.List `tfsdk:"video_codecs"`
}

type FfmpegConfigModel struct {
	Accel               types.String              `tfsdk:"accel"`
	AccelDecode         types.Bool                `tfsdk:"accel_decode"`
	AcceptedAudioCodecs types.List                `tfsdk:"accepted_audio_codecs"`
	AcceptedContainers  types.List                `tfsdk:"accepted_containers"`
	AcceptedVideoCodecs types.List                `tfsdk:"accepted_video_codecs"`
	Bframes             types.Int64               `tfsdk:"bframes"`
	CqMode              types.String              `tfsdk:"cq_mode"`
	Crf                 types.Int64               `tfsdk:"crf"`
	GopSize             types.Int64               `tfsdk:"gop_size"`
	MaxBitrate          types.String              `tfsdk:"max_bitrate"`
	PreferredHwDevice   types.String              `tfsdk:"preferred_hw_device"`
	Preset              types.String              `tfsdk:"preset"`
	Realtime            ffmpegRealtimeConfigModel `tfsdk:"realtime"`
	Refs                types.Int64               `tfsdk:"refs"`
	TargetAudioCodec    types.String              `tfsdk:"target_audio_codec"`
	TargetResolution    types.String              `tfsdk:"target_resolution"`
	TargetVideoCodec    types.String              `tfsdk:"target_video_codec"`
	TemporalAQ          types.Bool                `tfsdk:"temporal_aq"`
	Threads             types.Int64               `tfsdk:"threads"`
	Tonemap             types.String              `tfsdk:"tonemap"`
	Transcode           types.String              `tfsdk:"transcode"`
	TwoPass             types.Bool                `tfsdk:"two_pass"`
}

type fullsizeImageModel struct {
	Enabled     types.Bool   `tfsdk:"enabled"`
	Format      types.String `tfsdk:"format"`
	Progressive types.Bool   `tfsdk:"progressive"`
	Quality     types.Int64  `tfsdk:"quality"`
}

type imageModel struct {
	Format      types.String `tfsdk:"format"`
	Progressive types.Bool   `tfsdk:"progressive"`
	Quality     types.Int64  `tfsdk:"quality"`
	Size        types.Int64  `tfsdk:"size"`
}

type ImageConfigModel struct {
	ColorSpace      types.String       `tfsdk:"colorspace"`
	ExtractEmbedded types.Bool         `tfsdk:"extract_embedded"`
	Fullsize        fullsizeImageModel `tfsdk:"fullsize"`
	Preview         imageModel         `tfsdk:"preview"`
	Thumbnail       imageModel         `tfsdk:"thumbnail"`
}

type integrityChecksumJobConfigModel struct {
	CronExpression  types.String `tfsdk:"cron_expression"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	PercentageLimit types.Int64  `tfsdk:"percentage_limit"`
	TimeLimit       types.Int64  `tfsdk:"time_limit"`
}

type integrityJobConfigModel struct {
	CronExpression types.String `tfsdk:"cron_expression"`
	Enabled        types.Bool   `tfsdk:"enabled"`
}

type IntegrityChecksConfigModel struct {
	ChecksumFiles  integrityChecksumJobConfigModel `tfsdk:"checksum_files"`
	MissingFiles   integrityJobConfigModel         `tfsdk:"missing_files"`
	UntrackedFiles integrityJobConfigModel         `tfsdk:"untracked_files"`
}

type ServerConfigModel struct {
	ExternalDomain   types.String `tfsdk:"external_domain"`
	LoginPageMessage types.String `tfsdk:"login_page_message"`
	PublicUsers      types.Bool   `tfsdk:"public_users"`
}

type jobSettingConfigModel struct {
	Concurrency types.Int64 `tfsdk:"concurrency"`
}

type JobConfigModel struct {
	BackgroundTask      jobSettingConfigModel `tfsdk:"background_task"`
	Editor              jobSettingConfigModel `tfsdk:"editor"`
	FaceDetection       jobSettingConfigModel `tfsdk:"face_detection"`
	IntegrityCheck      jobSettingConfigModel `tfsdk:"integrity_check"`
	Library             jobSettingConfigModel `tfsdk:"library"`
	MetadataExtraction  jobSettingConfigModel `tfsdk:"metadata_extraction"`
	Migration           jobSettingConfigModel `tfsdk:"migration"`
	Notifications       jobSettingConfigModel `tfsdk:"notifications"`
	Ocr                 jobSettingConfigModel `tfsdk:"ocr"`
	Search              jobSettingConfigModel `tfsdk:"search"`
	Sidecar             jobSettingConfigModel `tfsdk:"sidecar"`
	SmartSearch         jobSettingConfigModel `tfsdk:"smart_search"`
	ThumbnailGeneration jobSettingConfigModel `tfsdk:"thumbnail_generation"`
	VideoConversion     jobSettingConfigModel `tfsdk:"video_conversion"`
	Workflow            jobSettingConfigModel `tfsdk:"workflow"`
}

type libraryScanConfigModel struct {
	CronExpression types.String `tfsdk:"cron_expression"`
	Enabled        types.Bool   `tfsdk:"enabled"`
}

type libraryWatchConfigModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type LibraryConfigModel struct {
	Scan  libraryScanConfigModel  `tfsdk:"scan"`
	Watch libraryWatchConfigModel `tfsdk:"watch"`
}

type LoggingConfigModel struct {
	Enabled types.Bool   `tfsdk:"enabled"`
	Level   types.String `tfsdk:"level"`
}

type machineLearningAvailabilityChecksConfigModel struct {
	Enabled  types.Bool  `tfsdk:"enabled"`
	Interval types.Int64 `tfsdk:"interval"`
	Timeout  types.Int64 `tfsdk:"timeout"`
}

type clipConfigModel struct {
	Enabled   types.Bool   `tfsdk:"enabled"`
	ModelName types.String `tfsdk:"model_name"`
}

type duplicateDetectionConfigModel struct {
	Enabled     types.Bool    `tfsdk:"enabled"`
	MaxDistance types.Float64 `tfsdk:"max_distance"`
}

type facialRecognitionConfigModel struct {
	Enabled     types.Bool    `tfsdk:"enabled"`
	MaxDistance types.Float64 `tfsdk:"max_distance"`
	MinFaces    types.Int64   `tfsdk:"min_faces"`
	MinScore    types.Float64 `tfsdk:"min_score"`
	ModelName   types.String  `tfsdk:"model_name"`
}

type ocrConfigModel struct {
	Enabled             types.Bool    `tfsdk:"enabled"`
	MaxResolution       types.Int64   `tfsdk:"max_resolution"`
	MinDetectionScore   types.Float64 `tfsdk:"min_detection_score"`
	MinRecognitionScore types.Float64 `tfsdk:"min_recognition_score"`
	ModelName           types.String  `tfsdk:"model_name"`
}

type MachineLearningConfigModel struct {
	AvailabilityChecks machineLearningAvailabilityChecksConfigModel `tfsdk:"availability_checks"`
	Clip               clipConfigModel                              `tfsdk:"clip"`
	DuplicateDetection duplicateDetectionConfigModel                `tfsdk:"duplicate_detection"`
	Enabled            types.Bool                                   `tfsdk:"enabled"`
	FacialRecognition  facialRecognitionConfigModel                 `tfsdk:"facial_recognition"`
	Ocr                ocrConfigModel                               `tfsdk:"ocr"`
	Urls               types.List                                   `tfsdk:"urls"`
}

type MapConfigModel struct {
	DarkStyle  types.String `tfsdk:"dark_style"`
	Enabled    types.Bool   `tfsdk:"enabled"`
	LightStyle types.String `tfsdk:"light_style"`
}

type facesConfigModel struct {
	Import types.Bool `tfsdk:"import"`
}

type MetadataConfigModel struct {
	Faces facesConfigModel `tfsdk:"faces"`
}

type NewVersionCheckConfigModel struct {
	Channel types.String `tfsdk:"channel"`
	Enabled types.Bool   `tfsdk:"enabled"`
}

type NightlyTasksConfigModel struct {
	ClusterNewFaces   types.Bool   `tfsdk:"cluster_new_faces"`
	DatabaseCleanup   types.Bool   `tfsdk:"database_cleanup"`
	GenerateMemories  types.Bool   `tfsdk:"generate_memories"`
	MissingThumbnails types.Bool   `tfsdk:"missing_thumbnails"`
	StartTime         types.String `tfsdk:"start_time"`
	SyncQuotaUsage    types.Bool   `tfsdk:"sync_quota_usage"`
}

type smtpTransportConfigModel struct {
	Host       types.String `tfsdk:"host"`
	IgnoreCert types.Bool   `tfsdk:"ignore_cert"`
	Password   types.String `tfsdk:"password"`
	Port       types.Int64  `tfsdk:"port"`
	Secure     types.Bool   `tfsdk:"secure"`
	Username   types.String `tfsdk:"username"`
}

type smtpConfigModel struct {
	Enabled   types.Bool               `tfsdk:"enabled"`
	From      types.String             `tfsdk:"from"`
	ReplyTo   types.String             `tfsdk:"reply_to"`
	Transport smtpTransportConfigModel `tfsdk:"transport"`
}

type NotificationsConfigModel struct {
	Smtp smtpConfigModel `tfsdk:"smtp"`
}

type OAuthConfigModel struct {
	AccountManagementUrl    types.String `tfsdk:"account_management_url"`
	AllowInsecureRequests   types.Bool   `tfsdk:"allow_insecure_requests"`
	AutoLaunch              types.Bool   `tfsdk:"auto_launch"`
	AutoRegister            types.Bool   `tfsdk:"auto_register"`
	ButtonText              types.String `tfsdk:"button_text"`
	ClientId                types.String `tfsdk:"client_id"`
	ClientSecret            types.String `tfsdk:"client_secret"`
	DefaultStorageQuota     types.Int64  `tfsdk:"default_storage_quota"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
	EndSessionEndpoint      types.String `tfsdk:"end_session_endpoint"`
	IssuerUrl               types.String `tfsdk:"issuer_url"`
	MobileOverrideEnabled   types.Bool   `tfsdk:"mobile_override_enabled"`
	MobileRedirectUri       types.String `tfsdk:"mobile_redirect_uri"`
	ProfileSigningAlgorithm types.String `tfsdk:"profile_signing_algorithm"`
	Prompt                  types.String `tfsdk:"prompt"`
	RoleClaim               types.String `tfsdk:"role_claim"`
	Scope                   types.String `tfsdk:"scope"`
	SigningAlgorithm        types.String `tfsdk:"signing_algorithm"`
	StorageLabelClaim       types.String `tfsdk:"storage_label_claim"`
	StorageQuotaClaim       types.String `tfsdk:"storage_quota_claim"`
	Timeout                 types.Int64  `tfsdk:"timeout"`
	TokenEndpointAuthMethod types.String `tfsdk:"token_endpoint_auth_method"`
}

type PasswordLoginConfigModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type ReverseGeocodingConfigModel struct {
	Enabled types.Bool `tfsdk:"enabled"`
}

type StorageTemplateConfigModel struct {
	Enabled                 types.Bool   `tfsdk:"enabled"`
	HashVerificationEnabled types.Bool   `tfsdk:"hash_verification_enabled"`
	Template                types.String `tfsdk:"template"`
}

type templateEmailsConfigModel struct {
	AlbumInviteTemplate types.String `tfsdk:"album_invite_template"`
	AlbumUpdateTemplate types.String `tfsdk:"album_update_template"`
	WelcomeTemplate     types.String `tfsdk:"welcome_template"`
}

type TemplatesConfigModel struct {
	Email templateEmailsConfigModel `tfsdk:"email"`
}

type ThemeConfigModel struct {
	CustomCss types.String `tfsdk:"custom_css"`
}

type TrashConfigModel struct {
	Days    types.Int64 `tfsdk:"days"`
	Enabled types.Bool  `tfsdk:"enabled"`
}

type UserConfigModel struct {
	DeleteDelay types.Int64 `tfsdk:"delete_delay"`
}

type configResourceModel struct {
	ID                     types.String                 `tfsdk:"id"`
	BackupsConfig          *BackupsConfigModel          `tfsdk:"backups"`
	FfmpegConfig           *FfmpegConfigModel           `tfsdk:"ffmpeg"`
	ImageConfig            *ImageConfigModel            `tfsdk:"image"`
	IntegrityChecksConfig  *IntegrityChecksConfigModel  `tfsdk:"integrity_checks"`
	JobConfig              *JobConfigModel              `tfsdk:"job"`
	LibraryConfig          *LibraryConfigModel          `tfsdk:"library"`
	LoggingConfig          *LoggingConfigModel          `tfsdk:"logging"`
	MachineLearningConfig  *MachineLearningConfigModel  `tfsdk:"machine_learning"`
	MapConfig              *MapConfigModel              `tfsdk:"map"`
	MetadataConfig         *MetadataConfigModel         `tfsdk:"metadata"`
	NewVersionCheckConfig  *NewVersionCheckConfigModel  `tfsdk:"new_version_check"`
	NightlyTasksConfig     *NightlyTasksConfigModel     `tfsdk:"nightly_tasks"`
	NotificationsConfig    *NotificationsConfigModel    `tfsdk:"notifications"`
	OauthConfig            *OAuthConfigModel            `tfsdk:"oauth"`
	PasswordLoginConfig    *PasswordLoginConfigModel    `tfsdk:"password_login"`
	ReverseGeocodingConfig *ReverseGeocodingConfigModel `tfsdk:"reverse_geocoding"`
	ServerConfig           *ServerConfigModel           `tfsdk:"server"`
	StorageTemplateConfig  *StorageTemplateConfigModel  `tfsdk:"storage_template"`
	TemplatesConfig        *TemplatesConfigModel        `tfsdk:"templates"`
	ThemeConfig            *ThemeConfigModel            `tfsdk:"theme"`
	TrashConfig            *TrashConfigModel            `tfsdk:"trash"`
	UserConfig             *UserConfigModel             `tfsdk:"user"`
}

func (r *configResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config"
}

func (r *configResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"backups": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Backup settings.",
				Attributes: map[string]schema.Attribute{
					"database": schema.SingleNestedAttribute{
						Required:    true,
						Description: "Database backup settings.",
						Attributes: map[string]schema.Attribute{
							"cron_expression": schema.StringAttribute{
								Description: "Set the scanning interval using the cron format.",
								Required:    true,
							},
							"enabled": schema.BoolAttribute{
								Description: "Enable database dumps.",
								Required:    true,
							},
							"keep_last_amount": schema.Int64Attribute{
								Description: "Number of previous dumps to keep.",
								Required:    true,
							},
						},
					},
				},
			},
			"ffmpeg": schema.SingleNestedAttribute{
				Required:    true,
				Description: "FFMpeg settings.",
				Attributes: map[string]schema.Attribute{
					"accel": schema.StringAttribute{
						Description: "The API that will interact with your device to accelerate transcoding. This setting is 'best effort': it will fallback to software transcoding on failure. VP9 may or may not work depending on your hardware.",
						Required:    true,
					},
					"accel_decode": schema.BoolAttribute{
						Description: "Enables end-to-end acceleration instead of only accelerating encoding. May not work on all videos.",
						Required:    true,
					},
					"accepted_audio_codecs": schema.ListAttribute{
						Description: "Select which audio codecs do not need to be transcoded. Only used for certain transcode policies.",
						ElementType: types.StringType,
						Required:    true,
					},
					"accepted_containers": schema.ListAttribute{
						Description: "Select which container formats do not need to be remuxed to MP4. Only used for certain transcode policies.",
						ElementType: types.StringType,
						Required:    true,
					},
					"accepted_video_codecs": schema.ListAttribute{
						Description: "Select which video codecs do not need to be transcoded. Only used for certain transcode policies.",
						ElementType: types.StringType,
						Required:    true,
					},
					"bframes": schema.Int64Attribute{
						Description: "Higher values improve compression efficiency, but slow down encoding. May not be compatible with hardware acceleration on older devices. 0 disables B-frames, while -1 sets this value automatically.",
						Required:    true,
					},
					"cq_mode": schema.StringAttribute{
						Description: "ICQ is better than CQP, but some hardware acceleration devices do not support this mode. Setting this option will prefer the specified mode when using quality-based encoding. Ignored by NVENC as it does not support ICQ.",
						Required:    true,
					},
					"crf": schema.Int64Attribute{
						Description: "Video quality level. Typical values are 23 for H.264, 28 for HEVC, 31 for VP9, and 35 for AV1. Lower is better, but produces larger files.",
						Required:    true,
					},
					"gop_size": schema.Int64Attribute{
						Description: "Sets the maximum frame distance between keyframes. Lower values worsen compression efficiency, but improve seek times and may improve quality in scenes with fast movement. 0 sets this value automatically.",
						Required:    true,
					},
					"max_bitrate": schema.StringAttribute{
						Description: "Setting a max bitrate can make file sizes more predictable at a minor cost to quality. At 720p, typical values are 2600 kbit/s for VP9 or HEVC, or 4500 kbit/s for H.264. Disabled if set to 0. When no unit is specified, k (for kbit/s) is assumed; therefore 5000, 5000k, and 5M (for Mbit/s) are equivalent.",
						Required:    true,
					},
					"preferred_hw_device": schema.StringAttribute{
						Description: "Applies only to VAAPI and QSV. Sets the dri node used for hardware transcoding.",
						Required:    true,
					},
					"preset": schema.StringAttribute{
						Description: "Compression speed. Slower presets produce smaller files, and increase quality when targeting a certain bitrate. VP9 ignores speeds above 'faster'.",
						Required:    true,
					},
					"realtime": schema.SingleNestedAttribute{
						Description: "Real-time Transcoding (experimental).",
						Required:    true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "If disabled, the server will refuse to start new real-time transcoding sessions.",
								Required:    true,
							},
							"resolutions": schema.ListAttribute{
								Description: "The resolutions offered for real-time transcoding. Higher resolutions may cause playback issues if the server cannot transcode them quickly enough.",
								ElementType: types.Int64Type,
								Required:    true,
							},
							"video_codecs": schema.ListAttribute{
								Description: "The video codecs offered for real-time transcoding. Clients will choose the best option they support during playback. AV1 is more efficient than HEVC, which is more efficient than H.264. When using hardware acceleration, only select the codecs the accelerator can encode. When using software transcoding, note that H.264 is faster than AV1, which is faster than HEVC.",
								ElementType: types.StringType,
								Required:    true,
							},
						},
					},
					"refs": schema.Int64Attribute{
						Description: "The number of frames to reference when compressing a given frame. Higher values improve compression efficiency, but slow down encoding. 0 sets this value automatically.",
						Required:    true,
					},
					"target_audio_codec": schema.StringAttribute{
						Description: "Opus is the highest quality option, but has lower compatibility with old devices or software.",
						Required:    true,
					},
					"target_resolution": schema.StringAttribute{
						Description: "Higher resolutions can preserve more detail but take longer to encode, have larger file sizes, and can reduce app responsiveness.",
						Required:    true,
					},
					"target_video_codec": schema.StringAttribute{
						Description: "VP9 has high efficiency and web compatibility, but takes longer to transcode. HEVC performs similarly, but has lower web compatibility. H.264 is widely compatible and quick to transcode, but produces much larger files. AV1 is the most efficient codec but lacks support on older devices.",
						Required:    true,
					},
					"temporal_aq": schema.BoolAttribute{
						Description: "Applies only to NVENC. Temporal Adaptive Quantisation increases quality of high-detail, low-motion scenes. May not be compatible with older devices.",
						Required:    true,
					},
					"threads": schema.Int64Attribute{
						Description: "Higher values lead to faster encoding, but leave less room for the server to process other tasks while active. This value should not be more than the number of CPU cores. Maximises utilisation if set to 0.",
						Required:    true,
					},
					"tonemap": schema.StringAttribute{
						Description: "Attempts to preserve the appearance of HDR videos when converted to SDR. Each algorithm makes different trade-offs for colour, detail and brightness. Hable preserves detail, Mobius preserves colour, and Reinhard preserves brightness.",
						Required:    true,
					},
					"transcode": schema.StringAttribute{
						Description: "Policy for when a video should be transcoded. HDR videos and videos with a pixel format other than YUV 4:2:0 will always be transcoded (except if transcoding is disabled).",
						Required:    true,
					},
					"two_pass": schema.BoolAttribute{
						Description: "Transcode in two passes to produce better encoded videos. When max bitrate is enabled (required for it to work with H.264 and HEVC), this mode uses a bitrate range based on the max bitrate and ignores CRF. For VP9, CRF can be used if max bitrate is disabled.",
						Required:    true,
					},
				},
			},
			"image": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Image settings.",
				Attributes: map[string]schema.Attribute{
					"colorspace": schema.StringAttribute{
						Description: "Use Display P3 for thumbnails. This better preserves the vibrance of images with wide colourspaces, but images may appear differently on old devices with an old browser version. sRGB images are kept as sRGB to avoid colour shifts.",
						Required:    true,
					},
					"extract_embedded": schema.BoolAttribute{
						Description: "Use embedded previews in RAW photos as the input to image processing and when available.",
						Required:    true,
					},
					"fullsize": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Generate full-size image for non-web-friendly formats.",
								Required:    true,
							},
							"format": schema.StringAttribute{
								Description: "Choose format of full-size image.",
								Required:    true,
							},
							"progressive": schema.BoolAttribute{
								Description: "Encode JPEG images progressively for gradual loading display.",
								Required:    true,
							},
							"quality": schema.Int64Attribute{
								Description: "Full-size image quality from 1-100. Higher is better, but produces larger files.",
								Required:    true,
							},
						},
					},
					"preview": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"format": schema.StringAttribute{
								Description: "Choose format of preview images.",
								Required:    true,
							},
							"progressive": schema.BoolAttribute{
								Description: "Encode JPEG images progressively for gradual loading display.",
								Required:    true,
							},
							"quality": schema.Int64Attribute{
								Description: "Preview image quality from 1-100. Higher is better, but produces larger files.",
								Required:    true,
							},
							"size": schema.Int64Attribute{
								Description: "Resolution of preview images.",
								Required:    true,
							},
						},
					},
					"thumbnail": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"format": schema.StringAttribute{
								Description: "Choose format of thumbnail images.",
								Required:    true,
							},
							"progressive": schema.BoolAttribute{
								Description: "Encode JPEG images progressively for gradual loading display.",
								Required:    true,
							},
							"quality": schema.Int64Attribute{
								Description: "Thumbnail image quality from 1-100. Higher is better, but produces larger files.",
								Required:    true,
							},
							"size": schema.Int64Attribute{
								Description: "Resolution of thumbnail images.",
								Required:    true,
							},
						},
					},
				},
			},
			"integrity_checks": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Integrity checks settings.",
				Attributes: map[string]schema.Attribute{
					"checksum_files": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"cron_expression": schema.StringAttribute{
								Description: "Set the scanning interval using the cron format.",
								Required:    true,
							},
							"enabled": schema.BoolAttribute{
								Description: "Enable the checksum check.",
								Required:    true,
							},
							"percentage_limit": schema.Int64Attribute{
								Description: "Configure the maximum percentage between 0.01 and 1 for how much the checksum check should run each interval.",
								Required:    true,
							},
							"time_limit": schema.Int64Attribute{
								Description: "Configure the maximum duration for which the checksum check should run each interval. (ms)",
								Required:    true,
							},
						},
					},
					"missing_files": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"cron_expression": schema.StringAttribute{
								Description: "Set the scanning interval using the cron format.",
								Required:    true,
							},
							"enabled": schema.BoolAttribute{
								Description: "Enable the missing files check.",
								Required:    true,
							},
						},
					},
					"untracked_files": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"cron_expression": schema.StringAttribute{
								Description: "Set the scanning interval using the cron format.",
								Required:    true,
							},
							"enabled": schema.BoolAttribute{
								Description: "Enable the untracked files check.",
								Required:    true,
							},
						},
					},
				},
			},
			"job": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Job settings.",
				Attributes: map[string]schema.Attribute{
					"background_task": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Background task concurrency.",
								Required:    true,
							},
						},
					},
					"editor": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Editor concurrency.",
								Required:    true,
							},
						},
					},
					"face_detection": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Face detection concurrency.",
								Required:    true,
							},
						},
					},
					"integrity_check": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Integrity checks concurrency.",
								Required:    true,
							},
						},
					},
					"library": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "External Libraries concurrency.",
								Required:    true,
							},
						},
					},
					"metadata_extraction": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Extract metadata concurrency.",
								Required:    true,
							},
						},
					},
					"migration": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Migration concurrency.",
								Required:    true,
							},
						},
					},
					"notifications": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Notifications concurrency.",
								Required:    true,
							},
						},
					},
					"ocr": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "OCR concurrency.",
								Required:    true,
							},
						},
					},
					"search": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Search concurrency.",
								Required:    true,
							},
						},
					},
					"sidecar": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Sidecar metadata concurrency.",
								Required:    true,
							},
						},
					},
					"smart_search": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Smart search concurrency.",
								Required:    true,
							},
						},
					},
					"thumbnail_generation": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Thumbnail generation concurrency.",
								Required:    true,
							},
						},
					},
					"video_conversion": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Video conversion concurrency.",
								Required:    true,
							},
						},
					},
					"workflow": schema.SingleNestedAttribute{
						Required: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Workflow concurrency.",
								Required:    true,
							},
						},
					},
				},
			},
			"library": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Library settings.",
				Attributes: map[string]schema.Attribute{
					"scan": schema.SingleNestedAttribute{
						Description: "Periodic scanning.",
						Required:    true,
						Attributes: map[string]schema.Attribute{
							"cron_expression": schema.StringAttribute{
								Description: "Set the scanning interval using the cron format.",
								Required:    true,
							},
							"enabled": schema.BoolAttribute{
								Description: "Enable periodic library scanning.",
								Required:    true,
							},
						},
					},
					"watch": schema.SingleNestedAttribute{
						Description: "Automatically watch for changed files.",
						Required:    true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Watch for changed files.",
								Required:    true,
							},
						},
					},
				},
			},
			"logging": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Logging settings.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: "Enable logging.",
						Required:    true,
					},
					"level": schema.StringAttribute{
						Description: "Sets the log level.",
						Required:    true,
					},
				},
			},
			"machine_learning": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Machine Learning settings.",
				Attributes: map[string]schema.Attribute{
					"availability_checks": schema.SingleNestedAttribute{
						Required:    true,
						Description: "Automatically detect and prefer available machine learning servers.",
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Set to true to enable.",
								Required:    true,
							},
							"interval": schema.Int64Attribute{
								Description: "Interval between checks, in milliseconds.",
								Required:    true,
							},
							"timeout": schema.Int64Attribute{
								Description: "Timeout for checks, in milliseconds.",
								Required:    true,
							},
						},
					},
					"clip": schema.SingleNestedAttribute{
						Required:    true,
						Description: "Use CLIP for smart search.",
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Enable CLIP.",
								Required:    true,
							},
							"model_name": schema.StringAttribute{
								Description: "Model to be used for CLIP.",
								Required:    true,
							},
						},
					},
					"duplicate_detection": schema.SingleNestedAttribute{
						Required:    true,
						Description: "Duplicate detection of images.",
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Enable duplicate detection.",
								Required:    true,
							},
							"max_distance": schema.Float64Attribute{
								Description: "Maximum distance between two images to consider them duplicates, ranging from 0.001-0.1. Higher values will detect more duplicates, but may result in false positives.",
								Required:    true,
							},
						},
					},
					"enabled": schema.BoolAttribute{
						Description: "Enable machine learning.",
						Required:    true,
					},
					"facial_recognition": schema.SingleNestedAttribute{
						Required:    true,
						Description: "Facial recognistion settings.",
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Enable facial recognition.",
								Required:    true,
							},
							"max_distance": schema.Float64Attribute{
								Description: "Maximum distance between two faces to be considered the same person, ranging from 0-2.",
								Required:    true,
							},
							"min_faces": schema.Int64Attribute{
								Description: "The minimum number of recognised faces for a person to be created.",
								Required:    true,
							},
							"min_score": schema.Float64Attribute{
								Description: "Minimum confidence score for a face to be detected from 0-1.",
								Required:    true,
							},
							"model_name": schema.StringAttribute{
								Description: "Facial recognition model to use.",
								Required:    true,
							},
						},
					},
					"ocr": schema.SingleNestedAttribute{
						Required:    true,
						Description: "OCR settings",
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Enable OCR.",
								Required:    true,
							},
							"max_resolution": schema.Int64Attribute{
								Description: "Previews above this resolution will be resized while preserving aspect ratio.",
								Required:    true,
							},
							"min_detection_score": schema.Float64Attribute{
								Description: "Minimum confidence score for text to be detected from 0-1.",
								Required:    true,
							},
							"min_recognition_score": schema.Float64Attribute{
								Description: "Minimum confidence score for detected text to be recognised from 0-1.",
								Required:    true,
							},
							"model_name": schema.StringAttribute{
								Description: "Model to use for OCR.",
								Required:    true,
							},
						},
					},
					"urls": schema.ListAttribute{
						Description: "The URL of the machine learning server. If more than one URL is provided, each server will be attempted one-at-a-time until one responds successfully, in order from first to last. Servers that don't respond will be temporarily ignored until they come back online.",
						ElementType: types.StringType,
						Required:    true,
					},
				},
			},
			"map": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Map settings.",
				Attributes: map[string]schema.Attribute{
					"dark_style": schema.StringAttribute{
						Description: "URL to a style.json map theme.",
						Required:    true,
					},
					"enabled": schema.BoolAttribute{
						Description: "The map feature relies on an external tile service (tiles.immich.cloud)",
						Required:    true,
					},
					"light_style": schema.StringAttribute{
						Description: "URL to a style.json map theme.",
						Required:    true,
					},
				},
			},
			"metadata": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Metadata settings.",
				Attributes: map[string]schema.Attribute{
					"faces": schema.SingleNestedAttribute{
						Required:    true,
						Description: "Face Metadata settings.",
						Attributes: map[string]schema.Attribute{
							"import": schema.BoolAttribute{
								Description: "Import faces from image EXIF data and sidecar files.",
								Required:    true,
							},
						},
					},
				},
			},
			"new_version_check": schema.SingleNestedAttribute{
				Required:    true,
				Description: "whatever",
				Attributes: map[string]schema.Attribute{
					"channel": schema.StringAttribute{
						Description: "Pick the release channel you want to get version announcements for.",
						Required:    true,
					},
					"enabled": schema.BoolAttribute{
						Description: "The version check feature relies on periodic communication with version.immich.cloud.",
						Required:    true,
					},
				},
			},
			"nightly_tasks": schema.SingleNestedAttribute{
				Required:    true,
				Description: "whatever",
				Attributes: map[string]schema.Attribute{
					"cluster_new_faces": schema.BoolAttribute{
						Description: "Run facial recognition on newly detected faces.",
						Required:    true,
					},
					"database_cleanup": schema.BoolAttribute{
						Description: "Clean up old, expired data from the database.",
						Required:    true,
					},
					"generate_memories": schema.BoolAttribute{
						Description: "Create new memories from assets.",
						Required:    true,
					},
					"missing_thumbnails": schema.BoolAttribute{
						Description: "Queue assets without thumbnails for thumbnail generation.",
						Required:    true,
					},
					"start_time": schema.StringAttribute{
						Description: "The time at which the server starts running the nightly tasks.",
						Required:    true,
					},
					"sync_quota_usage": schema.BoolAttribute{
						Description: "Update user storage quota, based on current usage.",
						Required:    true,
					},
				},
			},
			"notifications": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Notifications settings.",
				Attributes: map[string]schema.Attribute{
					"smtp": schema.SingleNestedAttribute{
						Required:    true,
						Description: "Email notifications settings.",
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Enable email notifications.",
								Required:    true,
							},
							"from": schema.StringAttribute{
								Description: "From address.",
								Required:    true,
							},
							"reply_to": schema.StringAttribute{
								Description: "Reply-To address.",
								Required:    true,
							},
							"transport": schema.SingleNestedAttribute{
								Required:    true,
								Description: "SMTP Transport settings.",
								Attributes: map[string]schema.Attribute{
									"host": schema.StringAttribute{
										Description: "Email server hostname.",
										Required:    true,
									},
									"ignore_cert": schema.BoolAttribute{
										Description: "Ignore TLS certificate validation errors (not recommended)",
										Required:    true,
									},
									"password": schema.StringAttribute{
										Description: "Password to use when authenticating with the email server.",
										Required:    true,
										Sensitive:   true,
									},
									"port": schema.Int64Attribute{
										Description: "Port of the email server (e.g 25, 465, or 587).",
										Required:    true,
									},
									"secure": schema.BoolAttribute{
										Description: "Use SMTPS (SMTP over TLS).",
										Required:    true,
									},
									"username": schema.StringAttribute{
										Description: "Username to use when authenticating with the email server.",
										Required:    true,
									},
								},
							},
						},
					},
				},
			},
			"oauth": schema.SingleNestedAttribute{
				Required:    true,
				Description: "OAuth settings.",
				Attributes: map[string]schema.Attribute{
					"account_management_url": schema.StringAttribute{
						Description: "Location in the external identity provider where a user can manage their settings or profile.",
						Required:    true,
					},
					"allow_insecure_requests": schema.BoolAttribute{
						Description: "WARNING: This disables TLS certificate validation for OAuth requests and may expose you to MITM attacks.",
						Required:    true,
					},
					"auto_launch": schema.BoolAttribute{
						Description: "Start the OAuth login flow automatically upon navigating to the login page.",
						Required:    true,
					},
					"auto_register": schema.BoolAttribute{
						Description: "Automatically register new users after signing in with OAuth.",
						Required:    true,
					},
					"button_text": schema.StringAttribute{
						Description: "Text on the login button.",
						Required:    true,
					},
					"client_id": schema.StringAttribute{
						Description: "Client Id",
						Required:    true,
					},
					"client_secret": schema.StringAttribute{
						Description: "Required for confidential client, or if PKCE (Proof Key for Code Exchange) is not supported for public client.",
						Required:    true,
						Sensitive:   true,
					},
					"default_storage_quota": schema.Int64Attribute{
						Description: "Quota in GiB to be used when no claim is provided.",
						Required:    true,
					},
					"enabled": schema.BoolAttribute{
						Description: "Enable OAuth.",
						Required:    true,
					},
					"end_session_endpoint": schema.StringAttribute{
						Description: "Redirect the user to this URI when they log out.",
						Required:    true,
					},
					"issuer_url": schema.StringAttribute{
						Description: "URL for OAuth issuer.",
						Required:    true,
					},
					"mobile_override_enabled": schema.BoolAttribute{
						Description: "Enable when OAuth provider does not allow a mobile URI, like 'app.immich:///oauth-callback'.",
						Required:    true,
					},
					"mobile_redirect_uri": schema.StringAttribute{
						Description: "Redirect URI for mobile.",
						Required:    true,
					},
					"profile_signing_algorithm": schema.StringAttribute{
						Description: "Algorithm for profile signing.",
						Required:    true,
					},
					"prompt": schema.StringAttribute{
						Description: "Prompt parameter (e.g. select_account, login, consent).",
						Required:    true,
					},
					"role_claim": schema.StringAttribute{
						Description: "Automatically grant admin access based on the presence of this claim. The claim may have either 'user' or 'admin'.",
						Required:    true,
					},
					"scope": schema.StringAttribute{
						Description: "Oauth scope.",
						Required:    true,
					},
					"signing_algorithm": schema.StringAttribute{
						Description: "Algorithm for signed response.",
						Required:    true,
					},
					"storage_label_claim": schema.StringAttribute{
						Description: "Automatically set the user's storage label to the value of this claim.",
						Required:    true,
					},
					"storage_quota_claim": schema.StringAttribute{
						Description: "Automatically set the user's storage quota to the value of this claim.",
						Required:    true,
					},
					"timeout": schema.Int64Attribute{
						Description: "Timeout for requests in milliseconds.",
						Required:    true,
					},
					"token_endpoint_auth_method": schema.StringAttribute{
						Description: "Endpoint login auth method.",
						Required:    true,
					},
				},
			},
			"password_login": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Password login settings.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: "Enable password login.",
						Required:    true,
					},
				},
			},
			"reverse_geocoding": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Reverse Geocoding settings.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: "Enable reverse geocoding.",
						Required:    true,
					},
				},
			},
			"server": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Server settings.",
				Attributes: map[string]schema.Attribute{
					"external_domain": schema.StringAttribute{
						Description: "Domain used for external links.",
						Required:    true,
					},
					"login_page_message": schema.StringAttribute{
						Description: "Welcome message displayed on login page.",
						Required:    true,
					},
					"public_users": schema.BoolAttribute{
						Description: "If true, all users are listed when adding users to shared albums.",
						Required:    true,
					},
				},
			},
			"storage_template": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Storage Template settings.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: "Enable storage template engine.",
						Required:    true,
					},
					"hash_verification_enabled": schema.BoolAttribute{
						Description: "Enable hash verification.",
						Required:    true,
					},
					"template": schema.StringAttribute{
						Description: "Storage template.",
						Required:    true,
					},
				},
			},
			"templates": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Templates settings.",
				Attributes: map[string]schema.Attribute{
					"email": schema.SingleNestedAttribute{
						Required:    true,
						Description: "whatever",
						Attributes: map[string]schema.Attribute{
							"album_invite_template": schema.StringAttribute{
								Required: true,
							},
							"album_update_template": schema.StringAttribute{
								Required: true,
							},
							"welcome_template": schema.StringAttribute{
								Required: true,
							},
						},
					},
				},
			},
			"theme": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Theme Settings.",
				Attributes: map[string]schema.Attribute{
					"custom_css": schema.StringAttribute{
						Description: "Cascading Style Sheets allow the design of Immich to be customised.",
						Required:    true,
					},
				},
			},
			"trash": schema.SingleNestedAttribute{
				Required:    true,
				Description: "Trash settings.",
				Attributes: map[string]schema.Attribute{
					"days": schema.Int64Attribute{
						Description: "Number of days to keep the assets in the bin before permanently removing them.",
						Required:    true,
					},
					"enabled": schema.BoolAttribute{
						Description: "Enable Trash features.",
						Required:    true,
					},
				},
			},
			"user": schema.SingleNestedAttribute{
				Required:    true,
				Description: "User settings.",
				Attributes: map[string]schema.Attribute{
					"delete_delay": schema.Int64Attribute{
						Description: "Number of days after removal to permanently delete a user's account and assets.",
						Required:    true,
					},
				},
			},
		},
	}
}

func (r *configResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state configResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	config, err := r.client.GetAdminConfig()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Config",
			"Could not read Immich Config:"+err.Error(),
		)
		return
	}

	if state.BackupsConfig == nil {
		state.BackupsConfig = &BackupsConfigModel{}
	}
	if state.FfmpegConfig == nil {
		state.FfmpegConfig = &FfmpegConfigModel{}
	}
	if state.ImageConfig == nil {
		state.ImageConfig = &ImageConfigModel{}
	}
	if state.IntegrityChecksConfig == nil {
		state.IntegrityChecksConfig = &IntegrityChecksConfigModel{}
	}
	if state.JobConfig == nil {
		state.JobConfig = &JobConfigModel{}
	}
	if state.LibraryConfig == nil {
		state.LibraryConfig = &LibraryConfigModel{}
	}
	if state.LoggingConfig == nil {
		state.LoggingConfig = &LoggingConfigModel{}
	}
	if state.MachineLearningConfig == nil {
		state.MachineLearningConfig = &MachineLearningConfigModel{}
	}
	if state.MapConfig == nil {
		state.MapConfig = &MapConfigModel{}
	}
	if state.MetadataConfig == nil {
		state.MetadataConfig = &MetadataConfigModel{}
	}
	if state.NewVersionCheckConfig == nil {
		state.NewVersionCheckConfig = &NewVersionCheckConfigModel{}
	}
	if state.NightlyTasksConfig == nil {
		state.NightlyTasksConfig = &NightlyTasksConfigModel{}
	}
	if state.NotificationsConfig == nil {
		state.NotificationsConfig = &NotificationsConfigModel{}
	}
	if state.OauthConfig == nil {
		state.OauthConfig = &OAuthConfigModel{}
	}
	if state.PasswordLoginConfig == nil {
		state.PasswordLoginConfig = &PasswordLoginConfigModel{}
	}
	if state.ReverseGeocodingConfig == nil {
		state.ReverseGeocodingConfig = &ReverseGeocodingConfigModel{}
	}
	if state.ServerConfig == nil {
		state.ServerConfig = &ServerConfigModel{}
	}
	if state.StorageTemplateConfig == nil {
		state.StorageTemplateConfig = &StorageTemplateConfigModel{}
	}
	if state.TemplatesConfig == nil {
		state.TemplatesConfig = &TemplatesConfigModel{}
	}
	if state.ThemeConfig == nil {
		state.ThemeConfig = &ThemeConfigModel{}
	}
	if state.TrashConfig == nil {
		state.TrashConfig = &TrashConfigModel{}
	}
	if state.UserConfig == nil {
		state.UserConfig = &UserConfigModel{}
	}

	state.BackupsConfig.Database.CronExpression = types.StringValue(config.Backup.Database.CronExpression)
	state.BackupsConfig.Database.Enabled = types.BoolValue(config.Backup.Database.Enabled)
	state.BackupsConfig.Database.KeepLastAmount = types.Int64Value(config.Backup.Database.KeepLastAmount)

	state.FfmpegConfig.Accel = types.StringValue(string(config.Ffmpeg.Accel))
	state.FfmpegConfig.AccelDecode = types.BoolValue(config.Ffmpeg.AccelDecode)

	audioCodecs, diags := types.ListValueFrom(ctx, types.StringType, config.Ffmpeg.AcceptedAudioCodecs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.FfmpegConfig.AcceptedAudioCodecs = audioCodecs

	containers, diags := types.ListValueFrom(ctx, types.StringType, config.Ffmpeg.AcceptedContainers)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.FfmpegConfig.AcceptedContainers = containers

	videoCodecs, diags := types.ListValueFrom(ctx, types.StringType, config.Ffmpeg.AcceptedVideoCodecs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.FfmpegConfig.AcceptedVideoCodecs = videoCodecs
	state.FfmpegConfig.Bframes = types.Int64Value(config.Ffmpeg.Bframes)
	state.FfmpegConfig.CqMode = types.StringValue(string(config.Ffmpeg.CqMode))
	state.FfmpegConfig.Crf = types.Int64Value(config.Ffmpeg.Crf)
	state.FfmpegConfig.GopSize = types.Int64Value(config.Ffmpeg.GopSize)
	state.FfmpegConfig.MaxBitrate = types.StringValue(config.Ffmpeg.MaxBitrate)
	state.FfmpegConfig.PreferredHwDevice = types.StringValue(config.Ffmpeg.PreferredHwDevice)
	state.FfmpegConfig.Preset = types.StringValue(config.Ffmpeg.Preset)
	state.FfmpegConfig.Realtime.Enabled = types.BoolValue(config.Ffmpeg.Realtime.Enabled)

	resolutions, diags := types.ListValueFrom(ctx, types.Int32Type, config.Ffmpeg.Realtime.Resolutions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.FfmpegConfig.Realtime.Resolutions = resolutions

	realtimeVideoCodecs, diags := types.ListValueFrom(ctx, types.StringType, config.Ffmpeg.Realtime.VideoCodecs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.FfmpegConfig.Realtime.VideoCodecs = realtimeVideoCodecs

	state.FfmpegConfig.Refs = types.Int64Value(config.Ffmpeg.Refs)
	state.FfmpegConfig.TargetAudioCodec = types.StringValue(string(config.Ffmpeg.TargetAudioCodec))
	state.FfmpegConfig.TargetResolution = types.StringValue(config.Ffmpeg.TargetResolution)
	state.FfmpegConfig.TargetVideoCodec = types.StringValue(string(config.Ffmpeg.TargetVideoCodec))
	state.FfmpegConfig.TemporalAQ = types.BoolValue(config.Ffmpeg.TemporalAQ)
	state.FfmpegConfig.Threads = types.Int64Value(config.Ffmpeg.Threads)
	state.FfmpegConfig.Tonemap = types.StringValue(string(config.Ffmpeg.Tonemap))
	state.FfmpegConfig.Transcode = types.StringValue(string(config.Ffmpeg.Transcode))
	state.FfmpegConfig.TwoPass = types.BoolValue(config.Ffmpeg.TwoPass)

	state.ImageConfig.ColorSpace = types.StringValue(string(config.Image.Colorspace))
	state.ImageConfig.ExtractEmbedded = types.BoolValue(config.Image.ExtractEmbedded)
	state.ImageConfig.Fullsize.Enabled = types.BoolValue(config.Image.Fullsize.Enabled)
	state.ImageConfig.Fullsize.Format = types.StringValue(string(config.Image.Fullsize.Format))
	state.ImageConfig.Fullsize.Progressive = types.BoolPointerValue(config.Image.Fullsize.Progressive)
	state.ImageConfig.Fullsize.Quality = types.Int64Value(config.Image.Fullsize.Quality)
	state.ImageConfig.Preview.Format = types.StringValue(string(config.Image.Preview.Format))
	state.ImageConfig.Preview.Progressive = types.BoolPointerValue(config.Image.Preview.Progressive)
	state.ImageConfig.Preview.Quality = types.Int64Value(config.Image.Preview.Quality)
	state.ImageConfig.Preview.Size = types.Int64Value(config.Image.Preview.Size)
	state.ImageConfig.Thumbnail.Format = types.StringValue(string(config.Image.Thumbnail.Format))
	state.ImageConfig.Thumbnail.Progressive = types.BoolPointerValue(config.Image.Thumbnail.Progressive)
	state.ImageConfig.Thumbnail.Quality = types.Int64Value(config.Image.Thumbnail.Quality)
	state.ImageConfig.Thumbnail.Size = types.Int64Value(config.Image.Thumbnail.Size)

	state.IntegrityChecksConfig.ChecksumFiles.CronExpression = types.StringValue(config.IntegrityChecks.ChecksumFiles.CronExpression)
	state.IntegrityChecksConfig.ChecksumFiles.Enabled = types.BoolValue(config.IntegrityChecks.ChecksumFiles.Enabled)
	state.IntegrityChecksConfig.ChecksumFiles.PercentageLimit = types.Int64Value(config.IntegrityChecks.ChecksumFiles.PercentageLimit)
	state.IntegrityChecksConfig.ChecksumFiles.TimeLimit = types.Int64Value(config.IntegrityChecks.ChecksumFiles.TimeLimit)
	state.IntegrityChecksConfig.MissingFiles.CronExpression = types.StringValue(config.IntegrityChecks.MissingFiles.CronExpression)
	state.IntegrityChecksConfig.MissingFiles.Enabled = types.BoolValue(config.IntegrityChecks.MissingFiles.Enabled)
	state.IntegrityChecksConfig.UntrackedFiles.CronExpression = types.StringValue(config.IntegrityChecks.UntrackedFiles.CronExpression)
	state.IntegrityChecksConfig.UntrackedFiles.Enabled = types.BoolValue(config.IntegrityChecks.UntrackedFiles.Enabled)

	state.JobConfig.BackgroundTask.Concurrency = types.Int64Value(config.Job.BackgroundTask.Concurrency)
	state.JobConfig.Editor.Concurrency = types.Int64Value(config.Job.Editor.Concurrency)
	state.JobConfig.FaceDetection.Concurrency = types.Int64Value(config.Job.FaceDetection.Concurrency)
	state.JobConfig.IntegrityCheck.Concurrency = types.Int64Value(config.Job.IntegrityCheck.Concurrency)
	state.JobConfig.Library.Concurrency = types.Int64Value(config.Job.Library.Concurrency)
	state.JobConfig.MetadataExtraction.Concurrency = types.Int64Value(config.Job.MetadataExtraction.Concurrency)
	state.JobConfig.Migration.Concurrency = types.Int64Value(config.Job.Migration.Concurrency)
	state.JobConfig.Notifications.Concurrency = types.Int64Value(config.Job.Notifications.Concurrency)
	state.JobConfig.Ocr.Concurrency = types.Int64Value(config.Job.Ocr.Concurrency)
	state.JobConfig.Search.Concurrency = types.Int64Value(config.Job.Search.Concurrency)
	state.JobConfig.Sidecar.Concurrency = types.Int64Value(config.Job.Sidecar.Concurrency)
	state.JobConfig.SmartSearch.Concurrency = types.Int64Value(config.Job.SmartSearch.Concurrency)
	state.JobConfig.ThumbnailGeneration.Concurrency = types.Int64Value(config.Job.ThumbnailGeneration.Concurrency)
	state.JobConfig.VideoConversion.Concurrency = types.Int64Value(config.Job.VideoConversion.Concurrency)
	state.JobConfig.Workflow.Concurrency = types.Int64Value(config.Job.Workflow.Concurrency)

	state.LibraryConfig.Scan.CronExpression = types.StringValue(config.Library.Scan.CronExpression)
	state.LibraryConfig.Scan.Enabled = types.BoolValue(config.Library.Scan.Enabled)
	state.LibraryConfig.Watch.Enabled = types.BoolValue(config.Library.Watch.Enabled)

	state.LoggingConfig.Enabled = types.BoolValue(config.Logging.Enabled)
	state.LoggingConfig.Level = types.StringValue(string(config.Logging.Level))

	state.MachineLearningConfig.AvailabilityChecks.Enabled = types.BoolValue(config.MachineLearning.AvailabilityChecks.Enabled)
	state.MachineLearningConfig.AvailabilityChecks.Interval = types.Int64Value(config.MachineLearning.AvailabilityChecks.Interval)
	state.MachineLearningConfig.AvailabilityChecks.Timeout = types.Int64Value(config.MachineLearning.AvailabilityChecks.Timeout)
	state.MachineLearningConfig.Clip.Enabled = types.BoolValue(config.MachineLearning.Clip.Enabled)
	state.MachineLearningConfig.Clip.ModelName = types.StringValue(config.MachineLearning.Clip.ModelName)
	state.MachineLearningConfig.DuplicateDetection.Enabled = types.BoolValue(config.MachineLearning.DuplicateDetection.Enabled)
	state.MachineLearningConfig.DuplicateDetection.MaxDistance = types.Float64Value(config.MachineLearning.DuplicateDetection.MaxDistance)
	state.MachineLearningConfig.Enabled = types.BoolValue(config.MachineLearning.Enabled)
	state.MachineLearningConfig.FacialRecognition.Enabled = types.BoolValue(config.MachineLearning.FacialRecognition.Enabled)
	state.MachineLearningConfig.FacialRecognition.MaxDistance = types.Float64Value(config.MachineLearning.FacialRecognition.MaxDistance)
	state.MachineLearningConfig.FacialRecognition.MinFaces = types.Int64Value(config.MachineLearning.FacialRecognition.MinFaces)
	state.MachineLearningConfig.FacialRecognition.MinScore = types.Float64Value(config.MachineLearning.FacialRecognition.MinScore)
	state.MachineLearningConfig.FacialRecognition.ModelName = types.StringValue(config.MachineLearning.FacialRecognition.ModelName)
	state.MachineLearningConfig.Ocr.Enabled = types.BoolValue(config.MachineLearning.Ocr.Enabled)
	state.MachineLearningConfig.Ocr.MaxResolution = types.Int64Value(config.MachineLearning.Ocr.MaxResolution)
	state.MachineLearningConfig.Ocr.MinDetectionScore = types.Float64Value(config.MachineLearning.Ocr.MinDetectionScore)
	state.MachineLearningConfig.Ocr.MinRecognitionScore = types.Float64Value(config.MachineLearning.Ocr.MinRecognitionScore)
	state.MachineLearningConfig.Ocr.ModelName = types.StringValue(config.MachineLearning.Ocr.ModelName)

	urls, diags := types.ListValueFrom(ctx, types.StringType, config.MachineLearning.Urls)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.MachineLearningConfig.Urls = urls

	state.MapConfig.DarkStyle = types.StringValue(string(config.Map.DarkStyle))
	state.MapConfig.Enabled = types.BoolValue(config.Map.Enabled)
	state.MapConfig.LightStyle = types.StringValue(string(config.Map.LightStyle))

	state.MetadataConfig.Faces.Import = types.BoolValue(config.Metadata.Faces.Import)

	state.NewVersionCheckConfig.Channel = types.StringValue(string(config.NewVersionCheck.Channel))
	state.NewVersionCheckConfig.Enabled = types.BoolValue(config.NewVersionCheck.Enabled)

	state.NightlyTasksConfig.ClusterNewFaces = types.BoolValue(config.NightlyTasks.ClusterNewFaces)
	state.NightlyTasksConfig.DatabaseCleanup = types.BoolValue(config.NightlyTasks.DatabaseCleanup)
	state.NightlyTasksConfig.GenerateMemories = types.BoolValue(config.NightlyTasks.GenerateMemories)
	state.NightlyTasksConfig.MissingThumbnails = types.BoolValue(config.NightlyTasks.MissingThumbnails)
	state.NightlyTasksConfig.StartTime = types.StringValue(config.NightlyTasks.StartTime)
	state.NightlyTasksConfig.SyncQuotaUsage = types.BoolValue(config.NightlyTasks.SyncQuotaUsage)

	state.NotificationsConfig.Smtp.Enabled = types.BoolValue(config.Notifications.Smtp.Enabled)
	state.NotificationsConfig.Smtp.From = types.StringValue(config.Notifications.Smtp.From)
	state.NotificationsConfig.Smtp.ReplyTo = types.StringValue(config.Notifications.Smtp.ReplyTo)
	state.NotificationsConfig.Smtp.Transport.Host = types.StringValue(config.Notifications.Smtp.Transport.Host)
	state.NotificationsConfig.Smtp.Transport.IgnoreCert = types.BoolValue(config.Notifications.Smtp.Transport.IgnoreCert)
	state.NotificationsConfig.Smtp.Transport.Password = types.StringValue(config.Notifications.Smtp.Transport.Password)
	state.NotificationsConfig.Smtp.Transport.Port = types.Int64Value(config.Notifications.Smtp.Transport.Port)
	state.NotificationsConfig.Smtp.Transport.Secure = types.BoolValue(config.Notifications.Smtp.Transport.Secure)
	state.NotificationsConfig.Smtp.Transport.Username = types.StringValue(config.Notifications.Smtp.Transport.Username)

	state.OauthConfig.AccountManagementUrl = types.StringPointerValue(config.Oauth.AccountManagementUrl)
	state.OauthConfig.AllowInsecureRequests = types.BoolValue(config.Oauth.AllowInsecureRequests)
	state.OauthConfig.AutoLaunch = types.BoolValue(config.Oauth.AutoLaunch)
	state.OauthConfig.AutoRegister = types.BoolValue(config.Oauth.AutoRegister)
	state.OauthConfig.ButtonText = types.StringValue(config.Oauth.ButtonText)
	state.OauthConfig.ClientId = types.StringValue(config.Oauth.ClientId)
	state.OauthConfig.ClientSecret = types.StringValue(config.Oauth.ClientSecret)
	state.OauthConfig.DefaultStorageQuota = types.Int64Value(config.Oauth.DefaultStorageQuota)
	state.OauthConfig.Enabled = types.BoolValue(config.Oauth.Enabled)
	state.OauthConfig.EndSessionEndpoint = types.StringValue(config.Oauth.EndSessionEndpoint)
	state.OauthConfig.IssuerUrl = types.StringValue(config.Oauth.IssuerUrl)
	state.OauthConfig.MobileOverrideEnabled = types.BoolValue(config.Oauth.MobileOverrideEnabled)
	state.OauthConfig.MobileRedirectUri = types.StringValue(config.Oauth.MobileRedirectUri)
	state.OauthConfig.ProfileSigningAlgorithm = types.StringValue(config.Oauth.ProfileSigningAlgorithm)
	state.OauthConfig.Prompt = types.StringValue(config.Oauth.Prompt)
	state.OauthConfig.RoleClaim = types.StringValue(config.Oauth.RoleClaim)
	state.OauthConfig.Scope = types.StringValue(config.Oauth.Scope)
	state.OauthConfig.SigningAlgorithm = types.StringValue(config.Oauth.SigningAlgorithm)
	state.OauthConfig.StorageLabelClaim = types.StringValue(config.Oauth.StorageLabelClaim)
	state.OauthConfig.StorageQuotaClaim = types.StringValue(config.Oauth.StorageQuotaClaim)
	state.OauthConfig.Timeout = types.Int64Value(config.Oauth.Timeout)
	state.OauthConfig.TokenEndpointAuthMethod = types.StringValue(string(config.Oauth.TokenEndpointAuthMethod))

	state.PasswordLoginConfig.Enabled = types.BoolValue(config.PasswordLogin.Enabled)

	state.ReverseGeocodingConfig.Enabled = types.BoolValue(config.ReverseGeocoding.Enabled)

	state.ServerConfig.ExternalDomain = types.StringValue(config.Server.ExternalDomain)
	state.ServerConfig.LoginPageMessage = types.StringValue(config.Server.LoginPageMessage)
	state.ServerConfig.PublicUsers = types.BoolValue(config.Server.PublicUsers)

	state.StorageTemplateConfig.Enabled = types.BoolValue(config.StorageTemplate.Enabled)
	state.StorageTemplateConfig.HashVerificationEnabled = types.BoolValue(config.StorageTemplate.HashVerificationEnabled)
	state.StorageTemplateConfig.Template = types.StringValue(config.StorageTemplate.Template)

	state.TemplatesConfig.Email.AlbumInviteTemplate = types.StringValue(config.Templates.Email.AlbumInviteTemplate)
	state.TemplatesConfig.Email.AlbumUpdateTemplate = types.StringValue(config.Templates.Email.AlbumUpdateTemplate)
	state.TemplatesConfig.Email.WelcomeTemplate = types.StringValue(config.Templates.Email.WelcomeTemplate)

	state.ThemeConfig.CustomCss = types.StringValue(config.Theme.CustomCss)

	state.TrashConfig.Days = types.Int64Value(config.Trash.Days)
	state.TrashConfig.Enabled = types.BoolValue(config.Trash.Enabled)
	state.UserConfig.DeleteDelay = types.Int64Value(config.User.DeleteDelay)

	state.ID = types.StringValue("singleton")

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *configResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan configResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var audioCodecs []immichclient.AudioCodec
	diags = plan.FfmpegConfig.AcceptedAudioCodecs.ElementsAs(ctx, &audioCodecs, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var containers []immichclient.VideoContainer
	diags = plan.FfmpegConfig.AcceptedContainers.ElementsAs(ctx, &containers, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var videoCodecs []immichclient.VideoCodec
	diags = plan.FfmpegConfig.AcceptedVideoCodecs.ElementsAs(ctx, &videoCodecs, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var resolutionsint32 []int32
	diags = plan.FfmpegConfig.Realtime.Resolutions.ElementsAs(ctx, &resolutionsint32, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var resolutions []immichclient.HlsVideoResolution
	for _, v := range resolutionsint32 {
		resolutions = append(resolutions, immichclient.HlsVideoResolution(v))
	}

	var realtimeVideoCodecs []immichclient.VideoCodec
	diags = plan.FfmpegConfig.Realtime.VideoCodecs.ElementsAs(ctx, &realtimeVideoCodecs, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var urls []string
	diags = plan.MachineLearningConfig.Urls.ElementsAs(ctx, &urls, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	configSettings := immichclient.AdminConfigDto{
		Backup: immichclient.AdminConfigBackupsDto{
			Database: immichclient.AdminConfigDatabaseBackupDto{
				CronExpression: plan.BackupsConfig.Database.CronExpression.ValueString(),
				Enabled:        plan.BackupsConfig.Database.Enabled.ValueBool(),
				KeepLastAmount: plan.BackupsConfig.Database.KeepLastAmount.ValueInt64(),
			},
		},
		Ffmpeg: immichclient.AdminConfigFFmpegDto{
			Accel:               immichclient.TranscodeHWAccel(plan.FfmpegConfig.Accel.ValueString()),
			AccelDecode:         plan.FfmpegConfig.AccelDecode.ValueBool(),
			AcceptedAudioCodecs: audioCodecs,
			AcceptedContainers:  containers,
			AcceptedVideoCodecs: videoCodecs,
			Bframes:             plan.FfmpegConfig.Bframes.ValueInt64(),
			CqMode:              immichclient.CQMode(plan.FfmpegConfig.CqMode.ValueString()),
			Crf:                 plan.FfmpegConfig.Crf.ValueInt64(),
			GopSize:             plan.FfmpegConfig.GopSize.ValueInt64(),
			MaxBitrate:          plan.FfmpegConfig.MaxBitrate.ValueString(),
			PreferredHwDevice:   plan.FfmpegConfig.PreferredHwDevice.ValueString(),
			Preset:              plan.FfmpegConfig.Preset.ValueString(),
			Realtime: immichclient.AdminConfigFFmpegRealtimeDto{
				Enabled:     plan.FfmpegConfig.Realtime.Enabled.ValueBool(),
				Resolutions: resolutions,
				VideoCodecs: realtimeVideoCodecs,
			},
			Refs:             plan.FfmpegConfig.Refs.ValueInt64(),
			TargetAudioCodec: immichclient.AudioCodec(plan.FfmpegConfig.TargetAudioCodec.ValueString()),
			TargetResolution: plan.FfmpegConfig.TargetResolution.ValueString(),
			TargetVideoCodec: immichclient.VideoCodec(plan.FfmpegConfig.TargetVideoCodec.ValueString()),
			TemporalAQ:       plan.FfmpegConfig.TemporalAQ.ValueBool(),
			Threads:          plan.FfmpegConfig.Threads.ValueInt64(),
			Tonemap:          immichclient.ToneMapping(plan.FfmpegConfig.Tonemap.ValueString()),
			Transcode:        immichclient.TranscodePolicy(plan.FfmpegConfig.Transcode.ValueString()),
			TwoPass:          plan.FfmpegConfig.TwoPass.ValueBool(),
		},
		Image: immichclient.AdminConfigImageDto{
			Colorspace:      immichclient.Colorspace(plan.ImageConfig.ColorSpace.ValueString()),
			ExtractEmbedded: plan.ImageConfig.ExtractEmbedded.ValueBool(),
			Fullsize: immichclient.AdminConfigGeneratedFullsizeImageDto{
				Enabled:     plan.ImageConfig.Fullsize.Enabled.ValueBool(),
				Format:      immichclient.ImageFormat(plan.ImageConfig.Fullsize.Format.ValueString()),
				Progressive: plan.ImageConfig.Fullsize.Progressive.ValueBoolPointer(),
				Quality:     plan.ImageConfig.Fullsize.Quality.ValueInt64(),
			},
			Preview: immichclient.AdminConfigGeneratedImageDto{
				Format:      immichclient.ImageFormat(plan.ImageConfig.Preview.Format.ValueString()),
				Progressive: plan.ImageConfig.Preview.Progressive.ValueBoolPointer(),
				Quality:     plan.ImageConfig.Preview.Quality.ValueInt64(),
				Size:        plan.ImageConfig.Preview.Size.ValueInt64(),
			},
			Thumbnail: immichclient.AdminConfigGeneratedImageDto{
				Format:      immichclient.ImageFormat(plan.ImageConfig.Thumbnail.Format.ValueString()),
				Progressive: plan.ImageConfig.Thumbnail.Progressive.ValueBoolPointer(),
				Quality:     plan.ImageConfig.Thumbnail.Quality.ValueInt64(),
				Size:        plan.ImageConfig.Thumbnail.Size.ValueInt64(),
			},
		},
		IntegrityChecks: immichclient.AdminConfigIntegrityChecksDto{
			ChecksumFiles: immichclient.AdminConfigIntegrityChecksumJobDto{
				CronExpression:  plan.IntegrityChecksConfig.ChecksumFiles.CronExpression.ValueString(),
				Enabled:         plan.IntegrityChecksConfig.ChecksumFiles.Enabled.ValueBool(),
				PercentageLimit: plan.IntegrityChecksConfig.ChecksumFiles.PercentageLimit.ValueInt64(),
				TimeLimit:       plan.IntegrityChecksConfig.ChecksumFiles.TimeLimit.ValueInt64(),
			},
			MissingFiles: immichclient.AdminConfigIntegrityJobDto{
				CronExpression: plan.IntegrityChecksConfig.MissingFiles.CronExpression.ValueString(),
				Enabled:        plan.IntegrityChecksConfig.MissingFiles.Enabled.ValueBool(),
			},
			UntrackedFiles: immichclient.AdminConfigIntegrityJobDto{
				CronExpression: plan.IntegrityChecksConfig.UntrackedFiles.CronExpression.ValueString(),
				Enabled:        plan.IntegrityChecksConfig.UntrackedFiles.Enabled.ValueBool(),
			},
		},
		Job: immichclient.AdminConfigJobDto{
			BackgroundTask: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.BackgroundTask.Concurrency.ValueInt64(),
			},
			Editor: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.Editor.Concurrency.ValueInt64(),
			},
			FaceDetection: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.FaceDetection.Concurrency.ValueInt64(),
			},
			IntegrityCheck: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.IntegrityCheck.Concurrency.ValueInt64(),
			},
			Library: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.Library.Concurrency.ValueInt64(),
			},
			MetadataExtraction: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.MetadataExtraction.Concurrency.ValueInt64(),
			},
			Migration: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.Migration.Concurrency.ValueInt64(),
			},
			Notifications: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.Notifications.Concurrency.ValueInt64(),
			},
			Ocr: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.Ocr.Concurrency.ValueInt64(),
			},
			Search: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.Search.Concurrency.ValueInt64(),
			},
			Sidecar: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.Sidecar.Concurrency.ValueInt64(),
			},
			SmartSearch: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.SmartSearch.Concurrency.ValueInt64(),
			},
			ThumbnailGeneration: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.ThumbnailGeneration.Concurrency.ValueInt64(),
			},
			VideoConversion: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.VideoConversion.Concurrency.ValueInt64(),
			},
			Workflow: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.Workflow.Concurrency.ValueInt64(),
			},
		},
		Library: immichclient.AdminConfigLibraryDto{
			Scan: immichclient.AdminConfigLibraryScanDto{
				CronExpression: plan.LibraryConfig.Scan.CronExpression.ValueString(),
				Enabled:        plan.LibraryConfig.Scan.Enabled.ValueBool(),
			},
			Watch: immichclient.AdminConfigLibraryWatchDto{
				Enabled: plan.LibraryConfig.Watch.Enabled.ValueBool(),
			},
		},
		Logging: immichclient.AdminConfigLoggingDto{
			Enabled: plan.LoggingConfig.Enabled.ValueBool(),
			Level:   immichclient.LogLevel(plan.LoggingConfig.Level.ValueString()),
		},
		MachineLearning: immichclient.AdminConfigMachineLearningDto{
			AvailabilityChecks: immichclient.AdminConfigMachineLearningAvailabilityChecksDto{
				Enabled:  plan.MachineLearningConfig.AvailabilityChecks.Enabled.ValueBool(),
				Interval: plan.MachineLearningConfig.AvailabilityChecks.Interval.ValueInt64(),
				Timeout:  plan.MachineLearningConfig.AvailabilityChecks.Timeout.ValueInt64(),
			},
			Clip: immichclient.AdminConfigClipDto{
				Enabled:   plan.MachineLearningConfig.Clip.Enabled.ValueBool(),
				ModelName: plan.MachineLearningConfig.Clip.ModelName.ValueString(),
			},
			DuplicateDetection: immichclient.AdminConfigDuplicateDetectionDto{
				Enabled:     plan.MachineLearningConfig.DuplicateDetection.Enabled.ValueBool(),
				MaxDistance: plan.MachineLearningConfig.DuplicateDetection.MaxDistance.ValueFloat64(),
			},
			Enabled: plan.MachineLearningConfig.Enabled.ValueBool(),
			FacialRecognition: immichclient.AdminConfigFacialRecognitionDto{
				Enabled:     plan.MachineLearningConfig.FacialRecognition.Enabled.ValueBool(),
				MaxDistance: plan.MachineLearningConfig.FacialRecognition.MaxDistance.ValueFloat64(),
				MinFaces:    plan.MachineLearningConfig.FacialRecognition.MinFaces.ValueInt64(),
				MinScore:    plan.MachineLearningConfig.FacialRecognition.MinScore.ValueFloat64(),
				ModelName:   plan.MachineLearningConfig.FacialRecognition.ModelName.ValueString(),
			},
			Ocr: immichclient.AdminConfigOcrDto{
				Enabled:             plan.MachineLearningConfig.Ocr.Enabled.ValueBool(),
				MaxResolution:       plan.MachineLearningConfig.Ocr.MaxResolution.ValueInt64(),
				MinDetectionScore:   plan.MachineLearningConfig.Ocr.MinDetectionScore.ValueFloat64(),
				MinRecognitionScore: plan.MachineLearningConfig.Ocr.MinRecognitionScore.ValueFloat64(),
				ModelName:           plan.MachineLearningConfig.Ocr.ModelName.ValueString(),
			},
			Urls: urls,
		},
		Map: immichclient.AdminConfigMapDto{
			DarkStyle:  plan.MapConfig.DarkStyle.ValueString(),
			Enabled:    plan.MapConfig.Enabled.ValueBool(),
			LightStyle: plan.MapConfig.LightStyle.ValueString(),
		},
		Metadata: immichclient.AdminConfigMetadataDto{
			Faces: immichclient.AdminConfigFacesDto{
				Import: plan.MetadataConfig.Faces.Import.ValueBool(),
			},
		},
		NewVersionCheck: immichclient.AdminConfigNewVersionCheckDto{
			Channel: immichclient.ReleaseChannel(plan.NewVersionCheckConfig.Channel.ValueString()),
			Enabled: plan.NewVersionCheckConfig.Enabled.ValueBool(),
		},
		NightlyTasks: immichclient.AdminConfigNightlyTasksDto{
			ClusterNewFaces:   plan.NightlyTasksConfig.ClusterNewFaces.ValueBool(),
			DatabaseCleanup:   plan.NightlyTasksConfig.DatabaseCleanup.ValueBool(),
			GenerateMemories:  plan.NightlyTasksConfig.GenerateMemories.ValueBool(),
			MissingThumbnails: plan.NightlyTasksConfig.MissingThumbnails.ValueBool(),
			StartTime:         plan.NightlyTasksConfig.StartTime.ValueString(),
			SyncQuotaUsage:    plan.NightlyTasksConfig.SyncQuotaUsage.ValueBool(),
		},
		Notifications: immichclient.AdminConfigNotificationsDto{
			Smtp: immichclient.AdminConfigSmtpDto{
				Enabled: plan.NotificationsConfig.Smtp.Enabled.ValueBool(),
				From:    plan.NotificationsConfig.Smtp.From.ValueString(),
				ReplyTo: plan.NotificationsConfig.Smtp.ReplyTo.ValueString(),
				Transport: immichclient.AdminConfigSmtpTransportDto{
					Host:       plan.NotificationsConfig.Smtp.Transport.Host.ValueString(),
					IgnoreCert: plan.NotificationsConfig.Smtp.Transport.IgnoreCert.ValueBool(),
					Password:   plan.NotificationsConfig.Smtp.Transport.Password.ValueString(),
					Port:       plan.NotificationsConfig.Smtp.Transport.Port.ValueInt64(),
					Secure:     plan.NotificationsConfig.Smtp.Transport.Secure.ValueBool(),
					Username:   plan.NotificationsConfig.Smtp.Transport.Username.ValueString(),
				},
			},
		},
		Oauth: immichclient.AdminConfigOAuthDto{
			AccountManagementUrl:    plan.OauthConfig.AccountManagementUrl.ValueStringPointer(),
			AllowInsecureRequests:   plan.OauthConfig.AllowInsecureRequests.ValueBool(),
			AutoLaunch:              plan.OauthConfig.AutoLaunch.ValueBool(),
			AutoRegister:            plan.OauthConfig.AutoRegister.ValueBool(),
			ButtonText:              plan.OauthConfig.ButtonText.ValueString(),
			ClientId:                plan.OauthConfig.ClientId.ValueString(),
			ClientSecret:            plan.OauthConfig.ClientSecret.ValueString(),
			DefaultStorageQuota:     plan.OauthConfig.DefaultStorageQuota.ValueInt64(),
			Enabled:                 plan.OauthConfig.Enabled.ValueBool(),
			EndSessionEndpoint:      plan.OauthConfig.EndSessionEndpoint.ValueString(),
			IssuerUrl:               plan.OauthConfig.IssuerUrl.ValueString(),
			MobileOverrideEnabled:   plan.OauthConfig.MobileOverrideEnabled.ValueBool(),
			MobileRedirectUri:       plan.OauthConfig.MobileRedirectUri.ValueString(),
			ProfileSigningAlgorithm: plan.OauthConfig.ProfileSigningAlgorithm.ValueString(),
			Prompt:                  plan.OauthConfig.Prompt.ValueString(),
			RoleClaim:               plan.OauthConfig.RoleClaim.ValueString(),
			Scope:                   plan.OauthConfig.Scope.ValueString(),
			SigningAlgorithm:        plan.OauthConfig.SigningAlgorithm.ValueString(),
			StorageLabelClaim:       plan.OauthConfig.StorageLabelClaim.ValueString(),
			StorageQuotaClaim:       plan.OauthConfig.StorageQuotaClaim.ValueString(),
			Timeout:                 plan.OauthConfig.Timeout.ValueInt64(),
			TokenEndpointAuthMethod: immichclient.OAuthTokenEndpointAuthMethod(plan.OauthConfig.TokenEndpointAuthMethod.ValueString()),
		},
		PasswordLogin: immichclient.AdminConfigPasswordLoginDto{
			Enabled: plan.PasswordLoginConfig.Enabled.ValueBool(),
		},
		ReverseGeocoding: immichclient.AdminConfigReverseGeocodingDto{
			Enabled: plan.ReverseGeocodingConfig.Enabled.ValueBool(),
		},
		Server: immichclient.AdminConfigServerDto{
			ExternalDomain:   plan.ServerConfig.ExternalDomain.ValueString(),
			LoginPageMessage: plan.ServerConfig.LoginPageMessage.ValueString(),
			PublicUsers:      plan.ServerConfig.PublicUsers.ValueBool(),
		},
		StorageTemplate: immichclient.AdminConfigStorageTemplateDto{
			Enabled:                 plan.StorageTemplateConfig.Enabled.ValueBool(),
			HashVerificationEnabled: plan.StorageTemplateConfig.HashVerificationEnabled.ValueBool(),
			Template:                plan.StorageTemplateConfig.Template.ValueString(),
		},
		Templates: immichclient.AdminConfigTemplatesDto{
			Email: immichclient.AdminConfigTemplateEmailsDto{
				AlbumInviteTemplate: plan.TemplatesConfig.Email.AlbumInviteTemplate.ValueString(),
				AlbumUpdateTemplate: plan.TemplatesConfig.Email.AlbumUpdateTemplate.ValueString(),
				WelcomeTemplate:     plan.TemplatesConfig.Email.WelcomeTemplate.ValueString(),
			},
		},
		Theme: immichclient.AdminConfigThemeDto{
			CustomCss: plan.ThemeConfig.CustomCss.ValueString(),
		},
		Trash: immichclient.AdminConfigTrashDto{
			Days:    plan.TrashConfig.Days.ValueInt64(),
			Enabled: plan.TrashConfig.Enabled.ValueBool(),
		},
		User: immichclient.AdminConfigUserDto{
			DeleteDelay: plan.UserConfig.DeleteDelay.ValueInt64(),
		},
	}

	_, err := r.client.UpdateAdminConfig(configSettings)

	plan.ID = types.StringValue("singleton")

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Immich Config",
			"Could not update order, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *configResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan configResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var audioCodecs []immichclient.AudioCodec
	diags = plan.FfmpegConfig.AcceptedAudioCodecs.ElementsAs(ctx, &audioCodecs, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var containers []immichclient.VideoContainer
	diags = plan.FfmpegConfig.AcceptedContainers.ElementsAs(ctx, &containers, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var videoCodecs []immichclient.VideoCodec
	diags = plan.FfmpegConfig.AcceptedVideoCodecs.ElementsAs(ctx, &videoCodecs, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var resolutionsint32 []int32
	diags = plan.FfmpegConfig.Realtime.Resolutions.ElementsAs(ctx, &resolutionsint32, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var resolutions []immichclient.HlsVideoResolution
	for _, v := range resolutionsint32 {
		resolutions = append(resolutions, immichclient.HlsVideoResolution(v))
	}

	var realtimeVideoCodecs []immichclient.VideoCodec
	diags = plan.FfmpegConfig.Realtime.VideoCodecs.ElementsAs(ctx, &realtimeVideoCodecs, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var urls []string
	diags = plan.MachineLearningConfig.Urls.ElementsAs(ctx, &urls, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	configSettings := immichclient.AdminConfigDto{
		Backup: immichclient.AdminConfigBackupsDto{
			Database: immichclient.AdminConfigDatabaseBackupDto{
				CronExpression: plan.BackupsConfig.Database.CronExpression.ValueString(),
				Enabled:        plan.BackupsConfig.Database.Enabled.ValueBool(),
				KeepLastAmount: plan.BackupsConfig.Database.KeepLastAmount.ValueInt64(),
			},
		},
		Ffmpeg: immichclient.AdminConfigFFmpegDto{
			Accel:               immichclient.TranscodeHWAccel(plan.FfmpegConfig.Accel.ValueString()),
			AccelDecode:         plan.FfmpegConfig.AccelDecode.ValueBool(),
			AcceptedAudioCodecs: audioCodecs,
			AcceptedContainers:  containers,
			AcceptedVideoCodecs: videoCodecs,
			Bframes:             plan.FfmpegConfig.Bframes.ValueInt64(),
			CqMode:              immichclient.CQMode(plan.FfmpegConfig.CqMode.ValueString()),
			Crf:                 plan.FfmpegConfig.Crf.ValueInt64(),
			GopSize:             plan.FfmpegConfig.GopSize.ValueInt64(),
			MaxBitrate:          plan.FfmpegConfig.MaxBitrate.ValueString(),
			PreferredHwDevice:   plan.FfmpegConfig.PreferredHwDevice.ValueString(),
			Preset:              plan.FfmpegConfig.Preset.ValueString(),
			Realtime: immichclient.AdminConfigFFmpegRealtimeDto{
				Enabled:     plan.FfmpegConfig.Realtime.Enabled.ValueBool(),
				Resolutions: resolutions,
				VideoCodecs: realtimeVideoCodecs,
			},
			Refs:             plan.FfmpegConfig.Refs.ValueInt64(),
			TargetAudioCodec: immichclient.AudioCodec(plan.FfmpegConfig.TargetAudioCodec.ValueString()),
			TargetResolution: plan.FfmpegConfig.TargetResolution.ValueString(),
			TargetVideoCodec: immichclient.VideoCodec(plan.FfmpegConfig.TargetVideoCodec.ValueString()),
			TemporalAQ:       plan.FfmpegConfig.TemporalAQ.ValueBool(),
			Threads:          plan.FfmpegConfig.Threads.ValueInt64(),
			Tonemap:          immichclient.ToneMapping(plan.FfmpegConfig.Tonemap.ValueString()),
			Transcode:        immichclient.TranscodePolicy(plan.FfmpegConfig.Transcode.ValueString()),
			TwoPass:          plan.FfmpegConfig.TwoPass.ValueBool(),
		},
		Image: immichclient.AdminConfigImageDto{
			Colorspace:      immichclient.Colorspace(plan.ImageConfig.ColorSpace.ValueString()),
			ExtractEmbedded: plan.ImageConfig.ExtractEmbedded.ValueBool(),
			Fullsize: immichclient.AdminConfigGeneratedFullsizeImageDto{
				Enabled:     plan.ImageConfig.Fullsize.Enabled.ValueBool(),
				Format:      immichclient.ImageFormat(plan.ImageConfig.Fullsize.Format.ValueString()),
				Progressive: plan.ImageConfig.Fullsize.Progressive.ValueBoolPointer(),
				Quality:     plan.ImageConfig.Fullsize.Quality.ValueInt64(),
			},
			Preview: immichclient.AdminConfigGeneratedImageDto{
				Format:      immichclient.ImageFormat(plan.ImageConfig.Preview.Format.ValueString()),
				Progressive: plan.ImageConfig.Preview.Progressive.ValueBoolPointer(),
				Quality:     plan.ImageConfig.Preview.Quality.ValueInt64(),
				Size:        plan.ImageConfig.Preview.Size.ValueInt64(),
			},
			Thumbnail: immichclient.AdminConfigGeneratedImageDto{
				Format:      immichclient.ImageFormat(plan.ImageConfig.Thumbnail.Format.ValueString()),
				Progressive: plan.ImageConfig.Thumbnail.Progressive.ValueBoolPointer(),
				Quality:     plan.ImageConfig.Thumbnail.Quality.ValueInt64(),
				Size:        plan.ImageConfig.Thumbnail.Size.ValueInt64(),
			},
		},
		IntegrityChecks: immichclient.AdminConfigIntegrityChecksDto{
			ChecksumFiles: immichclient.AdminConfigIntegrityChecksumJobDto{
				CronExpression:  plan.IntegrityChecksConfig.ChecksumFiles.CronExpression.ValueString(),
				Enabled:         plan.IntegrityChecksConfig.ChecksumFiles.Enabled.ValueBool(),
				PercentageLimit: plan.IntegrityChecksConfig.ChecksumFiles.PercentageLimit.ValueInt64(),
				TimeLimit:       plan.IntegrityChecksConfig.ChecksumFiles.TimeLimit.ValueInt64(),
			},
			MissingFiles: immichclient.AdminConfigIntegrityJobDto{
				CronExpression: plan.IntegrityChecksConfig.MissingFiles.CronExpression.ValueString(),
				Enabled:        plan.IntegrityChecksConfig.MissingFiles.Enabled.ValueBool(),
			},
			UntrackedFiles: immichclient.AdminConfigIntegrityJobDto{
				CronExpression: plan.IntegrityChecksConfig.UntrackedFiles.CronExpression.ValueString(),
				Enabled:        plan.IntegrityChecksConfig.UntrackedFiles.Enabled.ValueBool(),
			},
		},
		Job: immichclient.AdminConfigJobDto{
			BackgroundTask: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.BackgroundTask.Concurrency.ValueInt64(),
			},
			Editor: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.Editor.Concurrency.ValueInt64(),
			},
			FaceDetection: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.FaceDetection.Concurrency.ValueInt64(),
			},
			IntegrityCheck: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.IntegrityCheck.Concurrency.ValueInt64(),
			},
			Library: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.Library.Concurrency.ValueInt64(),
			},
			MetadataExtraction: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.MetadataExtraction.Concurrency.ValueInt64(),
			},
			Migration: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.Migration.Concurrency.ValueInt64(),
			},
			Notifications: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.Notifications.Concurrency.ValueInt64(),
			},
			Ocr: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.Ocr.Concurrency.ValueInt64(),
			},
			Search: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.Search.Concurrency.ValueInt64(),
			},
			Sidecar: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.Sidecar.Concurrency.ValueInt64(),
			},
			SmartSearch: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.SmartSearch.Concurrency.ValueInt64(),
			},
			ThumbnailGeneration: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.ThumbnailGeneration.Concurrency.ValueInt64(),
			},
			VideoConversion: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.VideoConversion.Concurrency.ValueInt64(),
			},
			Workflow: immichclient.AdminConfigJobSettingsDto{
				Concurrency: plan.JobConfig.Workflow.Concurrency.ValueInt64(),
			},
		},
		Library: immichclient.AdminConfigLibraryDto{
			Scan: immichclient.AdminConfigLibraryScanDto{
				CronExpression: plan.LibraryConfig.Scan.CronExpression.ValueString(),
				Enabled:        plan.LibraryConfig.Scan.Enabled.ValueBool(),
			},
			Watch: immichclient.AdminConfigLibraryWatchDto{
				Enabled: plan.LibraryConfig.Watch.Enabled.ValueBool(),
			},
		},
		Logging: immichclient.AdminConfigLoggingDto{
			Enabled: plan.LoggingConfig.Enabled.ValueBool(),
			Level:   immichclient.LogLevel(plan.LoggingConfig.Level.ValueString()),
		},
		MachineLearning: immichclient.AdminConfigMachineLearningDto{
			AvailabilityChecks: immichclient.AdminConfigMachineLearningAvailabilityChecksDto{
				Enabled:  plan.MachineLearningConfig.AvailabilityChecks.Enabled.ValueBool(),
				Interval: plan.MachineLearningConfig.AvailabilityChecks.Interval.ValueInt64(),
				Timeout:  plan.MachineLearningConfig.AvailabilityChecks.Timeout.ValueInt64(),
			},
			Clip: immichclient.AdminConfigClipDto{
				Enabled:   plan.MachineLearningConfig.Clip.Enabled.ValueBool(),
				ModelName: plan.MachineLearningConfig.Clip.ModelName.ValueString(),
			},
			DuplicateDetection: immichclient.AdminConfigDuplicateDetectionDto{
				Enabled:     plan.MachineLearningConfig.DuplicateDetection.Enabled.ValueBool(),
				MaxDistance: plan.MachineLearningConfig.DuplicateDetection.MaxDistance.ValueFloat64(),
			},
			Enabled: plan.MachineLearningConfig.Enabled.ValueBool(),
			FacialRecognition: immichclient.AdminConfigFacialRecognitionDto{
				Enabled:     plan.MachineLearningConfig.FacialRecognition.Enabled.ValueBool(),
				MaxDistance: plan.MachineLearningConfig.FacialRecognition.MaxDistance.ValueFloat64(),
				MinFaces:    plan.MachineLearningConfig.FacialRecognition.MinFaces.ValueInt64(),
				MinScore:    plan.MachineLearningConfig.FacialRecognition.MinScore.ValueFloat64(),
				ModelName:   plan.MachineLearningConfig.FacialRecognition.ModelName.ValueString(),
			},
			Ocr: immichclient.AdminConfigOcrDto{
				Enabled:             plan.MachineLearningConfig.Ocr.Enabled.ValueBool(),
				MaxResolution:       plan.MachineLearningConfig.Ocr.MaxResolution.ValueInt64(),
				MinDetectionScore:   plan.MachineLearningConfig.Ocr.MinDetectionScore.ValueFloat64(),
				MinRecognitionScore: plan.MachineLearningConfig.Ocr.MinRecognitionScore.ValueFloat64(),
				ModelName:           plan.MachineLearningConfig.Ocr.ModelName.ValueString(),
			},
			Urls: urls,
		},
		Map: immichclient.AdminConfigMapDto{
			DarkStyle:  plan.MapConfig.DarkStyle.ValueString(),
			Enabled:    plan.MapConfig.Enabled.ValueBool(),
			LightStyle: plan.MapConfig.LightStyle.ValueString(),
		},
		Metadata: immichclient.AdminConfigMetadataDto{
			Faces: immichclient.AdminConfigFacesDto{
				Import: plan.MetadataConfig.Faces.Import.ValueBool(),
			},
		},
		NewVersionCheck: immichclient.AdminConfigNewVersionCheckDto{
			Channel: immichclient.ReleaseChannel(plan.NewVersionCheckConfig.Channel.ValueString()),
			Enabled: plan.NewVersionCheckConfig.Enabled.ValueBool(),
		},
		NightlyTasks: immichclient.AdminConfigNightlyTasksDto{
			ClusterNewFaces:   plan.NightlyTasksConfig.ClusterNewFaces.ValueBool(),
			DatabaseCleanup:   plan.NightlyTasksConfig.DatabaseCleanup.ValueBool(),
			GenerateMemories:  plan.NightlyTasksConfig.GenerateMemories.ValueBool(),
			MissingThumbnails: plan.NightlyTasksConfig.MissingThumbnails.ValueBool(),
			StartTime:         plan.NightlyTasksConfig.StartTime.ValueString(),
			SyncQuotaUsage:    plan.NightlyTasksConfig.SyncQuotaUsage.ValueBool(),
		},
		Notifications: immichclient.AdminConfigNotificationsDto{
			Smtp: immichclient.AdminConfigSmtpDto{
				Enabled: plan.NotificationsConfig.Smtp.Enabled.ValueBool(),
				From:    plan.NotificationsConfig.Smtp.From.ValueString(),
				ReplyTo: plan.NotificationsConfig.Smtp.ReplyTo.ValueString(),
				Transport: immichclient.AdminConfigSmtpTransportDto{
					Host:       plan.NotificationsConfig.Smtp.Transport.Host.ValueString(),
					IgnoreCert: plan.NotificationsConfig.Smtp.Transport.IgnoreCert.ValueBool(),
					Password:   plan.NotificationsConfig.Smtp.Transport.Password.ValueString(),
					Port:       plan.NotificationsConfig.Smtp.Transport.Port.ValueInt64(),
					Secure:     plan.NotificationsConfig.Smtp.Transport.Secure.ValueBool(),
					Username:   plan.NotificationsConfig.Smtp.Transport.Username.ValueString(),
				},
			},
		},
		Oauth: immichclient.AdminConfigOAuthDto{
			AccountManagementUrl:    plan.OauthConfig.AccountManagementUrl.ValueStringPointer(),
			AllowInsecureRequests:   plan.OauthConfig.AllowInsecureRequests.ValueBool(),
			AutoLaunch:              plan.OauthConfig.AutoLaunch.ValueBool(),
			AutoRegister:            plan.OauthConfig.AutoRegister.ValueBool(),
			ButtonText:              plan.OauthConfig.ButtonText.ValueString(),
			ClientId:                plan.OauthConfig.ClientId.ValueString(),
			ClientSecret:            plan.OauthConfig.ClientSecret.ValueString(),
			DefaultStorageQuota:     plan.OauthConfig.DefaultStorageQuota.ValueInt64(),
			Enabled:                 plan.OauthConfig.Enabled.ValueBool(),
			EndSessionEndpoint:      plan.OauthConfig.EndSessionEndpoint.ValueString(),
			IssuerUrl:               plan.OauthConfig.IssuerUrl.ValueString(),
			MobileOverrideEnabled:   plan.OauthConfig.MobileOverrideEnabled.ValueBool(),
			MobileRedirectUri:       plan.OauthConfig.MobileRedirectUri.ValueString(),
			ProfileSigningAlgorithm: plan.OauthConfig.ProfileSigningAlgorithm.ValueString(),
			Prompt:                  plan.OauthConfig.Prompt.ValueString(),
			RoleClaim:               plan.OauthConfig.RoleClaim.ValueString(),
			Scope:                   plan.OauthConfig.Scope.ValueString(),
			SigningAlgorithm:        plan.OauthConfig.SigningAlgorithm.ValueString(),
			StorageLabelClaim:       plan.OauthConfig.StorageLabelClaim.ValueString(),
			StorageQuotaClaim:       plan.OauthConfig.StorageQuotaClaim.ValueString(),
			Timeout:                 plan.OauthConfig.Timeout.ValueInt64(),
			TokenEndpointAuthMethod: immichclient.OAuthTokenEndpointAuthMethod(plan.OauthConfig.TokenEndpointAuthMethod.ValueString()),
		},
		PasswordLogin: immichclient.AdminConfigPasswordLoginDto{
			Enabled: plan.PasswordLoginConfig.Enabled.ValueBool(),
		},
		ReverseGeocoding: immichclient.AdminConfigReverseGeocodingDto{
			Enabled: plan.ReverseGeocodingConfig.Enabled.ValueBool(),
		},
		Server: immichclient.AdminConfigServerDto{
			ExternalDomain:   plan.ServerConfig.ExternalDomain.ValueString(),
			LoginPageMessage: plan.ServerConfig.LoginPageMessage.ValueString(),
			PublicUsers:      plan.ServerConfig.PublicUsers.ValueBool(),
		},
		StorageTemplate: immichclient.AdminConfigStorageTemplateDto{
			Enabled:                 plan.StorageTemplateConfig.Enabled.ValueBool(),
			HashVerificationEnabled: plan.StorageTemplateConfig.HashVerificationEnabled.ValueBool(),
			Template:                plan.StorageTemplateConfig.Template.ValueString(),
		},
		Templates: immichclient.AdminConfigTemplatesDto{
			Email: immichclient.AdminConfigTemplateEmailsDto{
				AlbumInviteTemplate: plan.TemplatesConfig.Email.AlbumInviteTemplate.ValueString(),
				AlbumUpdateTemplate: plan.TemplatesConfig.Email.AlbumUpdateTemplate.ValueString(),
				WelcomeTemplate:     plan.TemplatesConfig.Email.WelcomeTemplate.ValueString(),
			},
		},
		Theme: immichclient.AdminConfigThemeDto{
			CustomCss: plan.ThemeConfig.CustomCss.ValueString(),
		},
		Trash: immichclient.AdminConfigTrashDto{
			Days:    plan.TrashConfig.Days.ValueInt64(),
			Enabled: plan.TrashConfig.Enabled.ValueBool(),
		},
		User: immichclient.AdminConfigUserDto{
			DeleteDelay: plan.UserConfig.DeleteDelay.ValueInt64(),
		},
	}

	_, err := r.client.UpdateAdminConfig(configSettings)

	plan.ID = types.StringValue("singleton")

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Immich Config",
			"Could not update order, unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *configResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}

func (r *configResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*immichclient.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *immichclient.Client, got %T. Please report this issue to the provider developer.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *configResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
