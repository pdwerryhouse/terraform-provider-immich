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

package provider

import (
	"context"
	"fmt"

	immichclient "codeberg.org/pdwerryhouse/immich-client-go/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &configDataSource{}
	_ datasource.DataSourceWithConfigure = &configDataSource{}
)

func NewConfigDataSource() datasource.DataSource {
	return &configDataSource{}
}

type configDataSource struct {
	client *immichclient.Client
}

type configDataSourceModel struct {
	BackupsConfig          backupsConfigModel          `tfsdk:"backups"`
	FfmpegConfig           ffmpegConfigModel           `tfsdk:"ffmpeg"`
	ImageConfig            imageConfigModel            `tfsdk:"image"`
	IntegrityChecksConfig  integrityChecksConfigModel  `tfsdk:"integrity_checks"`
	JobConfig              jobConfigModel              `tfsdk:"job"`
	LibraryConfig          libraryConfigModel          `tfsdk:"library"`
	LoggingConfig          loggingConfigModel          `tfsdk:"logging"`
	MachineLearningConfig  machineLearningConfigModel  `tfsdk:"machine_learning"`
	MapConfig              mapConfigModel              `tfsdk:"map"`
	MetadataConfig         metadataConfigModel         `tfsdk:"metadata"`
	NewVersionCheckConfig  newVersionCheckConfigModel  `tfsdk:"new_version_check"`
	NightlyTasksConfig     nightlyTasksConfigModel     `tfsdk:"nightly_tasks"`
	NotificationsConfig    notificationsConfigModel    `tfsdk:"notifications"`
	OauthConfig            oAuthConfigModel            `tfsdk:"oauth"`
	PasswordLoginConfig    passwordLoginConfigModel    `tfsdk:"password_login"`
	ReverseGeocodingConfig reverseGeocodingConfigModel `tfsdk:"reverse_geocoding"`
	ServerConfig           serverConfigModel           `tfsdk:"server"`
	StorageTemplateConfig  storageTemplateConfigModel  `tfsdk:"storage_template"`
	TemplatesConfig        templatesConfigModel        `tfsdk:"templates"`
	ThemeConfig            themeConfigModel            `tfsdk:"theme"`
	TrashConfig            trashConfigModel            `tfsdk:"trash"`
	UserConfig             userConfigModel             `tfsdk:"user"`
}

func (d *configDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_config"
}

func (d *configDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state configDataSourceModel

	config, err := d.client.GetAdminConfig()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Immich User",
			err.Error(),
		)
		return
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

	resolutions, diags := types.ListValueFrom(ctx, types.Int64Type, config.Ffmpeg.Realtime.Resolutions)
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

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *configDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*immichclient.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *immichclient.Client, got %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *configDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"backups": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Backup settings.",
				Attributes: map[string]schema.Attribute{
					"database": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Database backup settings.",
						Attributes: map[string]schema.Attribute{
							"cron_expression": schema.StringAttribute{
								Description: "Set the scanning interval using the cron format.",
								Computed:    true,
							},
							"enabled": schema.BoolAttribute{
								Description: "Enable database dumps.",
								Computed:    true,
							},
							"keep_last_amount": schema.Int64Attribute{
								Description: "Number of previous dumps to keep.",
								Computed:    true,
							},
						},
					},
				},
			},
			"ffmpeg": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "FFMpeg settings.",
				Attributes: map[string]schema.Attribute{
					"accel": schema.StringAttribute{
						Description: "The API that will interact with your device to accelerate transcoding. This setting is 'best effort': it will fallback to software transcoding on failure. VP9 may or may not work depending on your hardware.",
						Computed:    true,
					},
					"accel_decode": schema.BoolAttribute{
						Description: "Enables end-to-end acceleration instead of only accelerating encoding. May not work on all videos.",
						Computed:    true,
					},
					"accepted_audio_codecs": schema.ListAttribute{
						Description: "Select which audio codecs do not need to be transcoded. Only used for certain transcode policies.",
						ElementType: types.StringType,
						Computed:    true,
					},
					"accepted_containers": schema.ListAttribute{
						Description: "Select which container formats do not need to be remuxed to MP4. Only used for certain transcode policies.",
						ElementType: types.StringType,
						Computed:    true,
					},
					"accepted_video_codecs": schema.ListAttribute{
						Description: "Select which video codecs do not need to be transcoded. Only used for certain transcode policies.",
						ElementType: types.StringType,
						Computed:    true,
					},
					"bframes": schema.Int64Attribute{
						Description: "Higher values improve compression efficiency, but slow down encoding. May not be compatible with hardware acceleration on older devices. 0 disables B-frames, while -1 sets this value automatically.",
						Computed:    true,
					},
					"cq_mode": schema.StringAttribute{
						Description: "ICQ is better than CQP, but some hardware acceleration devices do not support this mode. Setting this option will prefer the specified mode when using quality-based encoding. Ignored by NVENC as it does not support ICQ.",
						Computed:    true,
					},
					"crf": schema.Int64Attribute{
						Description: "Video quality level. Typical values are 23 for H.264, 28 for HEVC, 31 for VP9, and 35 for AV1. Lower is better, but produces larger files.",
						Computed:    true,
					},
					"gop_size": schema.Int64Attribute{
						Description: "Sets the maximum frame distance between keyframes. Lower values worsen compression efficiency, but improve seek times and may improve quality in scenes with fast movement. 0 sets this value automatically.",
						Computed:    true,
					},
					"max_bitrate": schema.StringAttribute{
						Description: "Setting a max bitrate can make file sizes more predictable at a minor cost to quality. At 720p, typical values are 2600 kbit/s for VP9 or HEVC, or 4500 kbit/s for H.264. Disabled if set to 0. When no unit is specified, k (for kbit/s) is assumed; therefore 5000, 5000k, and 5M (for Mbit/s) are equivalent.",
						Computed:    true,
					},
					"preferred_hw_device": schema.StringAttribute{
						Description: "Applies only to VAAPI and QSV. Sets the dri node used for hardware transcoding.",
						Computed:    true,
					},
					"preset": schema.StringAttribute{
						Description: "Compression speed. Slower presets produce smaller files, and increase quality when targeting a certain bitrate. VP9 ignores speeds above 'faster'.",
						Computed:    true,
					},
					"realtime": schema.SingleNestedAttribute{
						Description: "Real-time Transcoding (experimental).",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "If disabled, the server will refuse to start new real-time transcoding sessions.",
								Computed:    true,
							},
							"resolutions": schema.ListAttribute{
								Description: "The resolutions offered for real-time transcoding. Higher resolutions may cause playback issues if the server cannot transcode them quickly enough.",
								ElementType: types.Int64Type,
								Computed:    true,
							},
							"video_codecs": schema.ListAttribute{
								Description: "The video codecs offered for real-time transcoding. Clients will choose the best option they support during playback. AV1 is more efficient than HEVC, which is more efficient than H.264. When using hardware acceleration, only select the codecs the accelerator can encode. When using software transcoding, note that H.264 is faster than AV1, which is faster than HEVC.",
								ElementType: types.StringType,
								Computed:    true,
							},
						},
					},
					"refs": schema.Int64Attribute{
						Description: "The number of frames to reference when compressing a given frame. Higher values improve compression efficiency, but slow down encoding. 0 sets this value automatically.",
						Computed:    true,
					},
					"target_audio_codec": schema.StringAttribute{
						Description: "Opus is the highest quality option, but has lower compatibility with old devices or software.",
						Computed:    true,
					},
					"target_resolution": schema.StringAttribute{
						Description: "Higher resolutions can preserve more detail but take longer to encode, have larger file sizes, and can reduce app responsiveness.",
						Computed:    true,
					},
					"target_video_codec": schema.StringAttribute{
						Description: "VP9 has high efficiency and web compatibility, but takes longer to transcode. HEVC performs similarly, but has lower web compatibility. H.264 is widely compatible and quick to transcode, but produces much larger files. AV1 is the most efficient codec but lacks support on older devices.",
						Computed:    true,
					},
					"temporal_aq": schema.BoolAttribute{
						Description: "Applies only to NVENC. Temporal Adaptive Quantisation increases quality of high-detail, low-motion scenes. May not be compatible with older devices.",
						Computed:    true,
					},
					"threads": schema.Int64Attribute{
						Description: "Higher values lead to faster encoding, but leave less room for the server to process other tasks while active. This value should not be more than the number of CPU cores. Maximises utilisation if set to 0.",
						Computed:    true,
					},
					"tonemap": schema.StringAttribute{
						Description: "Attempts to preserve the appearance of HDR videos when converted to SDR. Each algorithm makes different trade-offs for colour, detail and brightness. Hable preserves detail, Mobius preserves colour, and Reinhard preserves brightness.",
						Computed:    true,
					},
					"transcode": schema.StringAttribute{
						Description: "Policy for when a video should be transcoded. HDR videos and videos with a pixel format other than YUV 4:2:0 will always be transcoded (except if transcoding is disabled).",
						Computed:    true,
					},
					"two_pass": schema.BoolAttribute{
						Description: "Transcode in two passes to produce better encoded videos. When max bitrate is enabled (required for it to work with H.264 and HEVC), this mode uses a bitrate range based on the max bitrate and ignores CRF. For VP9, CRF can be used if max bitrate is disabled.",
						Computed:    true,
					},
				},
			},
			"image": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Image settings.",
				Attributes: map[string]schema.Attribute{
					"colorspace": schema.StringAttribute{
						Description: "Use Display P3 for thumbnails. This better preserves the vibrance of images with wide colourspaces, but images may appear differently on old devices with an old browser version. sRGB images are kept as sRGB to avoid colour shifts.",
						Computed:    true,
					},
					"extract_embedded": schema.BoolAttribute{
						Description: "Use embedded previews in RAW photos as the input to image processing and when available.",
						Computed:    true,
					},
					"fullsize": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Generate full-size image for non-web-friendly formats.",
								Computed:    true,
							},
							"format": schema.StringAttribute{
								Description: "Choose format of full-size image.",
								Computed:    true,
							},
							"progressive": schema.BoolAttribute{
								Description: "Encode JPEG images progressively for gradual loading display.",
								Computed:    true,
							},
							"quality": schema.Int64Attribute{
								Description: "Full-size image quality from 1-100. Higher is better, but produces larger files.",
								Computed:    true,
							},
						},
					},
					"preview": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"format": schema.StringAttribute{
								Description: "Choose format of preview images.",
								Computed:    true,
							},
							"progressive": schema.BoolAttribute{
								Description: "Encode JPEG images progressively for gradual loading display.",
								Computed:    true,
							},
							"quality": schema.Int64Attribute{
								Description: "Preview image quality from 1-100. Higher is better, but produces larger files.",
								Computed:    true,
							},
							"size": schema.Int64Attribute{
								Description: "Resolution of preview images.",
								Computed:    true,
							},
						},
					},
					"thumbnail": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"format": schema.StringAttribute{
								Description: "Choose format of thumbnail images.",
								Computed:    true,
							},
							"progressive": schema.BoolAttribute{
								Description: "Encode JPEG images progressively for gradual loading display.",
								Computed:    true,
							},
							"quality": schema.Int64Attribute{
								Description: "Thumbnail image quality from 1-100. Higher is better, but produces larger files.",
								Computed:    true,
							},
							"size": schema.Int64Attribute{
								Description: "Resolution of thumbnail images.",
								Computed:    true,
							},
						},
					},
				},
			},
			"integrity_checks": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Integrity checks settings.",
				Attributes: map[string]schema.Attribute{
					"checksum_files": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"cron_expression": schema.StringAttribute{
								Description: "Set the scanning interval using the cron format.",
								Computed:    true,
							},
							"enabled": schema.BoolAttribute{
								Description: "Enable the checksum check.",
								Computed:    true,
							},
							"percentage_limit": schema.Int64Attribute{
								Description: "Configure the maximum percentage between 0.01 and 1 for how much the checksum check should run each interval.",
								Computed:    true,
							},
							"time_limit": schema.Int64Attribute{
								Description: "Configure the maximum duration for which the checksum check should run each interval. (ms)",
								Computed:    true,
							},
						},
					},
					"missing_files": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"cron_expression": schema.StringAttribute{
								Description: "Set the scanning interval using the cron format.",
								Computed:    true,
							},
							"enabled": schema.BoolAttribute{
								Description: "Enable the missing files check.",
								Computed:    true,
							},
						},
					},
					"untracked_files": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"cron_expression": schema.StringAttribute{
								Description: "Set the scanning interval using the cron format.",
								Computed:    true,
							},
							"enabled": schema.BoolAttribute{
								Description: "Enable the untracked files check.",
								Computed:    true,
							},
						},
					},
				},
			},
			"job": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Job settings.",
				Attributes: map[string]schema.Attribute{
					"background_task": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Background task concurrency.",
								Computed:    true,
							},
						},
					},
					"editor": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Editor concurrency.",
								Computed:    true,
							},
						},
					},
					"face_detection": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Face detection concurrency.",
								Computed:    true,
							},
						},
					},
					"integrity_check": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Integrity checks concurrency.",
								Computed:    true,
							},
						},
					},
					"library": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "External Libraries concurrency.",
								Computed:    true,
							},
						},
					},
					"metadata_extraction": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Extract metadata concurrency.",
								Computed:    true,
							},
						},
					},
					"migration": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Migration concurrency.",
								Computed:    true,
							},
						},
					},
					"notifications": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Notifications concurrency.",
								Computed:    true,
							},
						},
					},
					"ocr": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "OCR concurrency.",
								Computed:    true,
							},
						},
					},
					"search": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Search concurrency.",
								Computed:    true,
							},
						},
					},
					"sidecar": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Sidecar metadata concurrency.",
								Computed:    true,
							},
						},
					},
					"smart_search": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Smart search concurrency.",
								Computed:    true,
							},
						},
					},
					"thumbnail_generation": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Thumbnail generation concurrency.",
								Computed:    true,
							},
						},
					},
					"video_conversion": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Video conversion concurrency.",
								Computed:    true,
							},
						},
					},
					"workflow": schema.SingleNestedAttribute{
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"concurrency": schema.Int64Attribute{
								Description: "Workflow concurrency.",
								Computed:    true,
							},
						},
					},
				},
			},
			"library": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Library settings.",
				Attributes: map[string]schema.Attribute{
					"scan": schema.SingleNestedAttribute{
						Description: "Periodic scanning.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"cron_expression": schema.StringAttribute{
								Description: "Set the scanning interval using the cron format.",
								Computed:    true,
							},
							"enabled": schema.BoolAttribute{
								Description: "Enable periodic library scanning.",
								Computed:    true,
							},
						},
					},
					"watch": schema.SingleNestedAttribute{
						Description: "Automatically watch for changed files.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Watch for changed files.",
								Computed:    true,
							},
						},
					},
				},
			},
			"logging": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Logging settings.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: "Enable logging.",
						Computed:    true,
					},
					"level": schema.StringAttribute{
						Description: "Sets the log level.",
						Computed:    true,
					},
				},
			},
			"machine_learning": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Machine Learning settings.",
				Attributes: map[string]schema.Attribute{
					"availability_checks": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Automatically detect and prefer available machine learning servers.",
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Set to true to enable.",
								Computed:    true,
							},
							"interval": schema.Int64Attribute{
								Description: "Interval between checks, in milliseconds.",
								Computed:    true,
							},
							"timeout": schema.Int64Attribute{
								Description: "Timeout for checks, in milliseconds.",
								Computed:    true,
							},
						},
					},
					"clip": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Use CLIP for smart search.",
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Enable CLIP.",
								Computed:    true,
							},
							"model_name": schema.StringAttribute{
								Description: "Model to be used for CLIP.",
								Computed:    true,
							},
						},
					},
					"duplicate_detection": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Duplicate detection of images.",
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Enable duplicate detection.",
								Computed:    true,
							},
							"max_distance": schema.Float64Attribute{
								Description: "Maximum distance between two images to consider them duplicates, ranging from 0.001-0.1. Higher values will detect more duplicates, but may result in false positives.",
								Computed:    true,
							},
						},
					},
					"enabled": schema.BoolAttribute{
						Description: "Enable machine learning.",
						Computed:    true,
					},
					"facial_recognition": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Facial recognistion settings.",
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Enable facial recognition.",
								Computed:    true,
							},
							"max_distance": schema.Float64Attribute{
								Description: "Maximum distance between two faces to be considered the same person, ranging from 0-2.",
								Computed:    true,
							},
							"min_faces": schema.Int64Attribute{
								Description: "The minimum number of recognised faces for a person to be created.",
								Computed:    true,
							},
							"min_score": schema.Float64Attribute{
								Description: "Minimum confidence score for a face to be detected from 0-1.",
								Computed:    true,
							},
							"model_name": schema.StringAttribute{
								Description: "Facial recognition model to use.",
								Computed:    true,
							},
						},
					},
					"ocr": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "OCR settings",
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Enable OCR.",
								Computed:    true,
							},
							"max_resolution": schema.Int64Attribute{
								Description: "Previews above this resolution will be resized while preserving aspect ratio.",
								Computed:    true,
							},
							"min_detection_score": schema.Float64Attribute{
								Description: "Minimum confidence score for text to be detected from 0-1.",
								Computed:    true,
							},
							"min_recognition_score": schema.Float64Attribute{
								Description: "Minimum confidence score for detected text to be recognised from 0-1.",
								Computed:    true,
							},
							"model_name": schema.StringAttribute{
								Description: "Model to use for OCR.",
								Computed:    true,
							},
						},
					},
					"urls": schema.ListAttribute{
						Description: "The URL of the machine learning server. If more than one URL is provided, each server will be attempted one-at-a-time until one responds successfully, in order from first to last. Servers that don't respond will be temporarily ignored until they come back online.",
						ElementType: types.StringType,
						Computed:    true,
					},
				},
			},
			"map": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Map settings.",
				Attributes: map[string]schema.Attribute{
					"dark_style": schema.StringAttribute{
						Description: "URL to a style.json map theme.",
						Computed:    true,
					},
					"enabled": schema.BoolAttribute{
						Description: "The map feature relies on an external tile service (tiles.immich.cloud)",
						Computed:    true,
					},
					"light_style": schema.StringAttribute{
						Description: "URL to a style.json map theme.",
						Computed:    true,
					},
				},
			},
			"metadata": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Metadata settings.",
				Attributes: map[string]schema.Attribute{
					"faces": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Face Metadata settings.",
						Attributes: map[string]schema.Attribute{
							"import": schema.BoolAttribute{
								Description: "Import faces from image EXIF data and sidecar files.",
								Computed:    true,
							},
						},
					},
				},
			},
			"new_version_check": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "whatever",
				Attributes: map[string]schema.Attribute{
					"channel": schema.StringAttribute{
						Description: "Pick the release channel you want to get version announcements for.",
						Computed:    true,
					},
					"enabled": schema.BoolAttribute{
						Description: "The version check feature relies on periodic communication with version.immich.cloud.",
						Computed:    true,
					},
				},
			},
			"nightly_tasks": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "whatever",
				Attributes: map[string]schema.Attribute{
					"cluster_new_faces": schema.BoolAttribute{
						Description: "Run facial recognition on newly detected faces.",
						Computed:    true,
					},
					"database_cleanup": schema.BoolAttribute{
						Description: "Clean up old, expired data from the database.",
						Computed:    true,
					},
					"generate_memories": schema.BoolAttribute{
						Description: "Create new memories from assets.",
						Computed:    true,
					},
					"missing_thumbnails": schema.BoolAttribute{
						Description: "Queue assets without thumbnails for thumbnail generation.",
						Computed:    true,
					},
					"start_time": schema.StringAttribute{
						Description: "The time at which the server starts running the nightly tasks.",
						Computed:    true,
					},
					"sync_quota_usage": schema.BoolAttribute{
						Description: "Update user storage quota, based on current usage.",
						Computed:    true,
					},
				},
			},
			"notifications": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Notifications settings.",
				Attributes: map[string]schema.Attribute{
					"smtp": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "Email notifications settings.",
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								Description: "Enable email notifications.",
								Computed:    true,
							},
							"from": schema.StringAttribute{
								Description: "From address.",
								Computed:    true,
							},
							"reply_to": schema.StringAttribute{
								Description: "Reply-To address.",
								Computed:    true,
							},
							"transport": schema.SingleNestedAttribute{
								Computed:    true,
								Description: "SMTP Transport settings.",
								Attributes: map[string]schema.Attribute{
									"host": schema.StringAttribute{
										Description: "Email server hostname.",
										Computed:    true,
									},
									"ignore_cert": schema.BoolAttribute{
										Description: "Ignore TLS certificate validation errors (not recommended)",
										Computed:    true,
									},
									"password": schema.StringAttribute{
										Description: "Password to use when authenticating with the email server.",
										Computed:    true,
									},
									"port": schema.Int64Attribute{
										Description: "Port of the email server (e.g 25, 465, or 587).",
										Computed:    true,
									},
									"secure": schema.BoolAttribute{
										Description: "Use SMTPS (SMTP over TLS).",
										Computed:    true,
									},
									"username": schema.StringAttribute{
										Description: "Username to use when authenticating with the email server.",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"oauth": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "OAuth settings.",
				Attributes: map[string]schema.Attribute{
					"account_management_url": schema.StringAttribute{
						Description: "Location in the external identity provider where a user can manage their settings or profile.",
						Computed:    true,
					},
					"allow_insecure_requests": schema.BoolAttribute{
						Description: "WARNING: This disables TLS certificate validation for OAuth requests and may expose you to MITM attacks.",
						Computed:    true,
					},
					"auto_launch": schema.BoolAttribute{
						Description: "Start the OAuth login flow automatically upon navigating to the login page.",
						Computed:    true,
					},
					"auto_register": schema.BoolAttribute{
						Description: "Automatically register new users after signing in with OAuth.",
						Computed:    true,
					},
					"button_text": schema.StringAttribute{
						Description: "Text on the login button.",
						Computed:    true,
					},
					"client_id": schema.StringAttribute{
						Description: "Client Id",
						Computed:    true,
					},
					"client_secret": schema.StringAttribute{
						Description: "Computed for confidential client, or if PKCE (Proof Key for Code Exchange) is not supported for public client.",
						Computed:    true,
					},
					"default_storage_quota": schema.Int64Attribute{
						Description: "Quota in GiB to be used when no claim is provided.",
						Computed:    true,
					},
					"enabled": schema.BoolAttribute{
						Description: "Enable OAuth.",
						Computed:    true,
					},
					"end_session_endpoint": schema.StringAttribute{
						Description: "Redirect the user to this URI when they log out.",
						Computed:    true,
					},
					"issuer_url": schema.StringAttribute{
						Description: "URL for OAuth issuer.",
						Computed:    true,
					},
					"mobile_override_enabled": schema.BoolAttribute{
						Description: "Enable when OAuth provider does not allow a mobile URI, like 'app.immich:///oauth-callback'.",
						Computed:    true,
					},
					"mobile_redirect_uri": schema.StringAttribute{
						Description: "Redirect URI for mobile.",
						Computed:    true,
					},
					"profile_signing_algorithm": schema.StringAttribute{
						Description: "Algorithm for profile signing.",
						Computed:    true,
					},
					"prompt": schema.StringAttribute{
						Description: "Prompt parameter (e.g. select_account, login, consent).",
						Computed:    true,
					},
					"role_claim": schema.StringAttribute{
						Description: "Automatically grant admin access based on the presence of this claim. The claim may have either 'user' or 'admin'.",
						Computed:    true,
					},
					"scope": schema.StringAttribute{
						Description: "Oauth scope.",
						Computed:    true,
					},
					"signing_algorithm": schema.StringAttribute{
						Description: "Algorithm for signed response.",
						Computed:    true,
					},
					"storage_label_claim": schema.StringAttribute{
						Description: "Automatically set the user's storage label to the value of this claim.",
						Computed:    true,
					},
					"storage_quota_claim": schema.StringAttribute{
						Description: "Automatically set the user's storage quota to the value of this claim.",
						Computed:    true,
					},
					"timeout": schema.Int64Attribute{
						Description: "Timeout for requests in milliseconds.",
						Computed:    true,
					},
					"token_endpoint_auth_method": schema.StringAttribute{
						Description: "Endpoint login auth method.",
						Computed:    true,
					},
				},
			},
			"password_login": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Password login settings.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: "Enable password login.",
						Computed:    true,
					},
				},
			},
			"reverse_geocoding": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Reverse Geocoding settings.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: "Enable reverse geocoding.",
						Computed:    true,
					},
				},
			},
			"server": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Server settings.",
				Attributes: map[string]schema.Attribute{
					"external_domain": schema.StringAttribute{
						Description: "Domain used for external links.",
						Computed:    true,
					},
					"login_page_message": schema.StringAttribute{
						Description: "Welcome message displayed on login page.",
						Computed:    true,
					},
					"public_users": schema.BoolAttribute{
						Description: "If true, all users are listed when adding users to shared albums.",
						Computed:    true,
					},
				},
			},
			"storage_template": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Storage Template settings.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Description: "Enable storage template engine.",
						Computed:    true,
					},
					"hash_verification_enabled": schema.BoolAttribute{
						Description: "Enable hash verification.",
						Computed:    true,
					},
					"template": schema.StringAttribute{
						Description: "Storage template.",
						Computed:    true,
					},
				},
			},
			"templates": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Templates settings.",
				Attributes: map[string]schema.Attribute{
					"email": schema.SingleNestedAttribute{
						Computed:    true,
						Description: "whatever",
						Attributes: map[string]schema.Attribute{
							"album_invite_template": schema.StringAttribute{
								Computed: true,
							},
							"album_update_template": schema.StringAttribute{
								Computed: true,
							},
							"welcome_template": schema.StringAttribute{
								Computed: true,
							},
						},
					},
				},
			},
			"theme": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Theme Settings.",
				Attributes: map[string]schema.Attribute{
					"custom_css": schema.StringAttribute{
						Description: "Cascading Style Sheets allow the design of Immich to be customised.",
						Computed:    true,
					},
				},
			},
			"trash": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Trash settings.",
				Attributes: map[string]schema.Attribute{
					"days": schema.Int64Attribute{
						Description: "Number of days to keep the assets in the bin before permanently removing them.",
						Computed:    true,
					},
					"enabled": schema.BoolAttribute{
						Description: "Enable Trash features.",
						Computed:    true,
					},
				},
			},
			"user": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "User settings.",
				Attributes: map[string]schema.Attribute{
					"delete_delay": schema.Int64Attribute{
						Description: "Number of days after removal to permanently delete a user's account and assets.",
						Computed:    true,
					},
				},
			},
		},
	}

}
