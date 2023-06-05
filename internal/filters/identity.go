package filters

import (
	"github.com/lusis/statusthing/internal/serrors"
	"github.com/lusis/statusthing/internal/validation"
)

// WithEmailAddress provides a custom [v1.User] email address
func WithEmailAddress(email string) FilterOption {
	return func(f *Filters) error {
		if !validation.ValidString(email) {
			return serrors.NewError("email", serrors.ErrEmptyString)
		}
		if f.emailaddress != nil {
			return serrors.NewError("email", serrors.ErrAlreadySet)
		}
		f.emailaddress = &email
		return nil
	}
}

// EmailAddress returns the configured [v1.User] email address
func (f *Filters) EmailAddress() string {
	f.l.RLock()
	defer f.l.RUnlock()
	if f.emailaddress == nil {
		return ""
	}
	return *f.emailaddress
}

// WithFirstName provides a custom [v1.User] firstname
func WithFirstName(fname string) FilterOption {
	return func(f *Filters) error {
		if !validation.ValidString(fname) {
			return serrors.NewError("fname", serrors.ErrEmptyString)
		}
		if f.firstname != nil {
			return serrors.NewError("email", serrors.ErrAlreadySet)
		}
		f.firstname = &fname
		return nil
	}
}

// FirstName returns the configured [v1.User] first name
func (f *Filters) FirstName() string {
	f.l.RLock()
	defer f.l.RUnlock()
	if f.firstname == nil {
		return ""
	}
	return *f.firstname
}

// WithLastName provides a custom [v1.User] last name
func WithLastName(lname string) FilterOption {
	return func(f *Filters) error {
		if !validation.ValidString(lname) {
			return serrors.NewError("lname", serrors.ErrEmptyString)
		}
		if f.lastname != nil {
			return serrors.NewError("lname", serrors.ErrAlreadySet)
		}
		f.lastname = &lname
		return nil
	}
}

// LastName returns the configured [v1.User] last name
func (f *Filters) LastName() string {
	f.l.RLock()
	defer f.l.RUnlock()
	if f.lastname == nil {
		return ""
	}
	return *f.lastname
}
