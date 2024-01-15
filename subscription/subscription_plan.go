package subscription

const (
	Basic        = "Basic"        // 0 usd limited
	Essential    = "Essential"    // 7 usd/year 9 usd/month
	Premium      = "Premium"      // 15 usd/year 18 usd/month
	Professional = "Professional" // 24 usd/year 27 usd/month
	Elite        = "Elite"        // 33 usd/year 36 usd/month
)

type SubsriptionPlan struct {
	Name        string
	Description string
	Price       float64

	//Meetings
	MeetingScheduleProfilesLimit          int
	AvailabilityProfilesLimit             int
	MeetingsScheduleCustomFieldsAvailable bool
	MeetingsAvailabilityOverrides         bool
	TeamsAvailable                        bool
	TeamMembersLimit                      int
	TeamsLimit                            int
	AIAvatarsAvailable                    bool
	AIAvatarsLimit                        int
	MeetingsInsightsAvailable             bool
	GoogleCalendarSupport                 bool
	GoogleMeetSupport                     bool
	ZoomSupport                           bool
	PayedMeetingsSupport                  bool
	ContactsLimit                         int
	MeetingsLimit                         int
	ChatsAvailable                        bool
	TeamChatsAvailable                    bool
	EmailNotificationsAvailable           bool
	SmsNotificationsAvailable             bool
	WorksflowsAvailable                   bool
	WorkflowsLimit                        int

	//Todos
	TodoCategoriesLimit int
	TodoItemsLimit      int

	//Habits
	Limit                   int
	HabitTemplatesLimit     int
	HabitTemplatesAvailable bool

	//General
	AiAssistantAvailable bool
}
