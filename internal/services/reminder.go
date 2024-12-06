package services

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/DukeRupert/rr/internal/email"
	"github.com/DukeRupert/rr/internal/orderspace"
	"github.com/go-co-op/gocron/v2"
)

type ReminderScheduler struct {
	scheduler gocron.Scheduler
}

func NewReminderScheduler(db *sql.DB, orderClient *orderspace.Client, emailClient *email.Client) (*ReminderScheduler, error) {
	s, err := gocron.NewScheduler()
	if err != nil {
		return nil, fmt.Errorf("creating scheduler: %w", err)
	}

	_, err = s.NewJob(
		gocron.WeeklyJob(
			1,

			gocron.NewWeekdays(time.Friday),
			gocron.NewAtTimes(gocron.NewAtTime(9, 0, 0)),
		),
		gocron.NewTask(
			func() error {
				return sendOrderReminders(db, orderClient, emailClient)
			},
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating reminder job: %w", err)
	}

	return &ReminderScheduler{scheduler: s}, nil
}

func (rs *ReminderScheduler) Start() {
	rs.scheduler.Start()
}

func (rs *ReminderScheduler) Shutdown() error {
	return rs.scheduler.Shutdown()
}

func sendOrderReminders(db *sql.DB, orderClient *orderspace.Client, emailClient *email.Client) error {
	// Get customers with orders in last 6 weeks
	sixWeeksAgo := time.Now().AddDate(0, 0, -42)
	params := &orderspace.CustomerListParams{
		UpdatedSince: &sixWeeksAgo,
	}

	resp, err := orderClient.ListCustomers(params)
	if err != nil {
		return fmt.Errorf("fetching customers: %w", err)
	}

	for _, customer := range resp.Customers {
		// Check if customer has opted out
		var notifyDays bool
		err := db.QueryRow(`
            INSERT INTO customer_notifications (customer_id, email_notify_days)
            VALUES (?, true)
            ON CONFLICT (customer_id) DO UPDATE SET updated_at = CURRENT_TIMESTAMP
            RETURNING email_notify_days
        `, customer.ID).Scan(&notifyDays)
		if err != nil {
			return fmt.Errorf("checking notification preference: %w", err)
		}

		if !notifyDays {
			continue
		}

		// Send reminder email
		reminderEmail := email.Email{
			From:     "info@rockabillyroasting.com",
			To:       customer.EmailAddresses.Orders,
			Subject:  "Reminder: Place Your Order by Monday",
			HtmlBody: generateReminderEmailHTML(customer.CompanyName),
			TextBody: generateReminderEmailText(customer.CompanyName),
		}

		_, err = emailClient.SendEmail(reminderEmail)
		if err != nil {
			return fmt.Errorf("sending reminder to %s: %w", customer.ID, err)
		}
	}

	return nil
}

func generateReminderEmailHTML(companyName string) string {
	return fmt.Sprintf(`
        <html>
            <body>
                <h2>Hey there, %s!</h2>
                <p>Just a friendly reminder from your coffee crew at Rockabilly Roasting over here in Washington State.</p>
                <p>To keep your coffee delivery running smooth as a '57 Chevy, we kindly ask that you place your order by Monday. This helps us make sure your beans arrive right on schedule the following week.</p>
                <p>Need anything else? Just hit reply - we're always happy to help!</p>
                <p>Keep rockin',<br>
                The Rockabilly Roasting Team</p>
            </body>
        </html>
    `, companyName)
}

func generateReminderEmailText(companyName string) string {
	return fmt.Sprintf(`Hey there, %s!

Just a friendly reminder from your coffee crew at Rockabilly Roasting over here in Washington State.

To keep your coffee delivery running smooth as a '57 Chevy, we kindly ask that you place your order by Monday. This helps us make sure your beans arrive right on schedule the following week.

Need anything else? Just hit reply - we're always happy to help!

Keep rockin',
The Rockabilly Roasting Team`, companyName)
}

func PreviewOrderReminders(db *sql.DB, orderClient *orderspace.Client, emailClient *email.Client) error {
	sixWeeksAgo := time.Now().AddDate(0, 0, -42)
	params := &orderspace.CustomerListParams{
		UpdatedSince: &sixWeeksAgo,
	}

	resp, err := orderClient.ListCustomers(params)
	if err != nil {
		return fmt.Errorf("fetching customers: %w", err)
	}

	var activeCustomers []string
	for _, customer := range resp.Customers {
		var notifyDays bool
		err := db.QueryRow(`
            SELECT COALESCE(
                (SELECT email_notify_days FROM customer_notifications WHERE customer_id = ?),
                true
            )
        `, customer.ID).Scan(&notifyDays)
		if err != nil {
			return fmt.Errorf("checking notification preference: %w", err)
		}

		if notifyDays {
			activeCustomers = append(activeCustomers, fmt.Sprintf("%s (%s)", customer.CompanyName, customer.EmailAddresses.Orders))
		}
	}

	// Send preview email
	previewEmail := email.Email{
		From:     "info@rockabillyroasting.com",
		To:       "logan@fireflysoftware.dev",
		Subject:  fmt.Sprintf("Order Reminder Preview - %d Customers", len(activeCustomers)),
		HtmlBody: generatePreviewEmailHTML(activeCustomers),
		TextBody: generatePreviewEmailText(activeCustomers),
	}

	_, err = emailClient.SendEmail(previewEmail)
	return err
}

func generatePreviewEmailHTML(customers []string) string {
	customerList := strings.Join(customers, "<br>")
	return fmt.Sprintf(`
        <html>
            <body>
                <h2>Rockabilly Roasting Reminder Preview</h2>
                <p>Hey there! Here's who's getting our friendly Monday order reminder this week:</p>
                <p><strong>%d customers on the list:</strong></p>
                <p>%s</p>
                <hr>
                <p><em>These customers will receive our standard reminder about placing orders by Monday for next week's delivery.</em></p>
            </body>
        </html>
    `, len(customers), customerList)
}

func generatePreviewEmailText(customers []string) string {
	customerList := strings.Join(customers, "\n")
	return fmt.Sprintf(`Rockabilly Roasting Reminder Preview

Hey there! Here's who's getting our friendly Monday order reminder this week:

%d customers on the list:

%s

These customers will receive our standard reminder about placing orders by Monday for next week's delivery.`,
		len(customers), customerList)
}
