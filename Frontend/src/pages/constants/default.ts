import type { StudentAcademicBackground, StudentContactDetails, StudentProfile, StudentSettings } from "../../services/profile.service";


export const defaultStudentProfile: StudentProfile = {
  gender: "male",
  nationality: "",
  nationalId: "",
  dateOfBirth: "",
  passportNumber: "",
};

export const defaultContact: StudentContactDetails = {
  email: "",
  phoneNumber: "",
  preferredContactMethod: "email",
  addressLine: "",
  city: "",
  country: "",
  zipCode: "",
};

export const defaultAcademic: StudentAcademicBackground = {
  qualification: "",
  institutionName: "",
  graduationYear: new Date().getFullYear(),
  gpaScore: "",
  languageTestType: "",
  languageTestScore: "",
  certificateFilePath: "",
};

export const defaultSettings: StudentSettings = {
  languagePreference: "en",
  receiveEmailNotifications: false,
  receiveSMSNotifications: false,
  emergencyName: "",
  emergencyRelationship: "",
  emergencyPhoneNumber: "",
  emergencyEmail: "",
  accountDeleted: false,
};
