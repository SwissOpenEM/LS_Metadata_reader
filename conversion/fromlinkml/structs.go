package oscem
import "LS_reader/conversion/basetypes"


/*
 * basetypes.String type, used as the base for type-narrowing.
 */
type Any struct {
}

/*
 * A range constructed from min and max
 */
type Range struct {
	/*
	 * Minimal value of a given dataset property
	 */
	Minimal basetypes.Float64 `json:"minimal"`
	/*
	 * Maximal value of a given dataset property
	 */
	Maximal basetypes.Float64 `json:"maximal"`
}

/*
 * A series of numbers constructed from min, max, and increment
 */
type Series struct {
	/*
	 * parent types
	 */
	Range
	/*
	 * Increment between elements of a series
	 */
	Increment basetypes.Float64 `json:"increment"`
	/*
	 * Minimal value of a given dataset property
	 */
	Minimal basetypes.Float64 `json:"minimal"`
	/*
	 * Maximal value of a given dataset property
	 */
	Maximal basetypes.Float64 `json:"maximal"`
}

/*
 * size of a 2D image (in basetypes.Integer units)
 */
type ImageSize struct {
	/*
	 * The height of a given item - unit depends on item
	 */
	Height basetypes.Int `json:"height"`
	/*
	 * The width of a given item - unit depends on item
	 */
	Width basetypes.Int `json:"width"`
}

/*
 * an axis-aligned 2D bounding box (float units)
 */
type BoundingBox2D struct {
	/*
	 * minimum x
	 */
	XMin basetypes.Float64 `json:"x_min"`
	/*
	 * maximum x
	 */
	XMax basetypes.Float64 `json:"x_max"`
	/*
	 * minimum y
	 */
	YMin basetypes.Float64 `json:"y_min"`
	/*
	 * maximum y
	 */
	YMax basetypes.Float64 `json:"y_max"`
}

/*
 * if a value has a unit, it should be given as a unit value pair.
 */
type QuantityValue struct {
	/*
	 * the unit of a given value
	 */
	Unit basetypes.String `json:"unit"`
	/*
	 * the value of a field with a unit
	 */
	Value basetypes.Float64 `json:"value"`
}

/*
 * A set of parameteres describing the data acquisition
 */
type Acquisition struct {
	/*
	 * Target defocus set, min and max values in µm.
	 */
	NominalDefocus Range `json:"nominal_defocus"`
	/*
	 * Machine estimated defocus, min and max values in µm. Has a tendency to be off.
	 */
	CalibratedDefocus Range `json:"calibrated_defocus"`
	/*
	 * Magnification level as indicated by the instrument, no unit
	 */
	NominalMagnification basetypes.Int `json:"nominal_magnification"`
	/*
	 * Calculated magnification, no unit
	 */
	CalibratedMagnification basetypes.Int `json:"calibrated_magnification"`
	/*
	 * Speciman holder model
	 */
	Holder basetypes.String `json:"holder"`
	/*
	 * Type of cryogen used in the holder - if the holder is cooled seperately
	 */
	HolderCryogen basetypes.String `json:"holder_cryogen"`
	/*
	 * Temperature during data collection, in K with min and max values.
	 */
	Temperature Range `json:"temperature"`
	/*
	 * Software used for instrument control,
	 */
	MicroscopeSoftware basetypes.String `json:"microscope_software"`
	/*
	 * Make and model of the detector used
	 */
	Detector basetypes.String `json:"detector"`
	/*
	 * Operating mode of the detector
	 */
	DetectorMode basetypes.String `json:"detector_mode"`
	/*
	 * Average dose per image/movie/tilt - given in electrons per square Angstrom
	 */
	DosePerMovie basetypes.Float64 `json:"dose_per_movie"`
	/*
	 * Wether an energy filter was used and its specifics.
	 */
	EnergyFilter EnergyFilter `json:"energy_filter"`
	/*
	 * The size of the image in pixels, height and width given.
	 */
	ImageSize ImageSize `json:"image_size"`
	/*
	 * Time and date of the data acquisition
	 */
	DateTime basetypes.String `json:"date_time"`
	/*
	 * Time of data acquisition per movie/tilt - in s
	 */
	ExposureTime basetypes.Float64 `json:"exposure_time"`
	/*
	 * Cryogen used in cooling the instrument and sample, usually nitrogen
	 */
	Cryogen basetypes.String `json:"cryogen"`
	/*
	 * Number of frames that on average constitute a full movie, can be a bit hard to define for some detectors
	 */
	FramesPerMovie basetypes.Int `json:"frames_per_movie"`
	/*
	 * Number of grids imaged for this project - here with qualifier during this data acquisition
	 */
	GridsImaged basetypes.Int `json:"grids_imaged"`
	/*
	 * Number of images generated total for this data collection - might need a qualifier for tilt series to determine whether full series or individual tilts are counted
	 */
	ImagesGenerated basetypes.Int `json:"images_generated"`
	/*
	 * Level of binning on the images applied during data collection
	 */
	BinningCamera basetypes.Float64 `json:"binning_camera"`
	/*
	 * Pixel size, in Angstrom
	 */
	PixelSize basetypes.Float64 `json:"pixel_size"`
	/*
	 * basetypes.String type of special optics, such as a phaseplate
	 */
	SpecialistOptics SpecialistOptics `json:"specialist_optics"`
	/*
	 * Movement of the beam above the sample for data collection purposes that does not require movement of the stage. Given in mrad.
	 */
	Beamshift BoundingBox2D `json:"beamshift"`
	/*
	 * Another way to move the beam above the sample for data collection purposes that does not require movement of the stage. Given in mrad.
	 */
	Beamtilt BoundingBox2D `json:"beamtilt"`
	/*
	 * Movement of the Beam below the image in order to shift the image on the detector. Given in µm.
	 */
	Imageshift BoundingBox2D `json:"imageshift"`
	/*
	 * Number of Beamtilt groups present in this dataset - for optimized processing split dataset basetypes.Into groups of same tilt angle. Despite its name Beamshift is often used to achive this result.
	 */
	Beamtiltgroups basetypes.Int `json:"beamtiltgroups"`
	/*
	 * Whether and how you have to flip or rotate the gainref in order to align with your acquired images
	 */
	GainrefFlipRotate basetypes.String `json:"gainref_flip_rotate"`
}

/*
 * A device used to filter for electrons with specific energy.
 */
type EnergyFilter struct {
	/*
	 * whether a specific instrument was used during data acquisition
	 */
	Used basetypes.Bool `json:"used"`
	/*
	 * Make and model of a specilized device
	 */
	Model basetypes.String `json:"model"`
	/*
	 * Width of the energy filter used.
	 */
	WidthEnergyFilter basetypes.Float64 `json:"width_energy_filter"`
}

/*
 * Optional optics used to correct for instrument limitations.
 */
type SpecialistOptics struct {
	/*
	 * Phaseplate is a special optics device that can be used to enhance contrast
	 */
	Phaseplate Phaseplate `json:"phaseplate"`
	/*
	 * Specialist device to correct for spherical aberration of the microscope lenses
	 */
	SphericalAberrationCorrector SphericalAberrationCorrector `json:"spherical_aberration_corrector"`
	/*
	 * Specialist device to correct for chromatic aberration of the microscope lenses
	 */
	ChromaticAberrationCorrector ChromaticAberrationCorrector `json:"chromatic_aberration_corrector"`
}

/*
 * Used to modulate the phase of the electron wave.
 */
type Phaseplate struct {
	/*
	 * whether a specific instrument was used during data acquisition
	 */
	Used basetypes.Bool `json:"used"`
	/*
	 * Type of phaseplate
	 */
	InstrumentType basetypes.String `json:"instrument_type"`
}

/*
 * Special device used to correct instrument inherent spherical aberration.
 */
type SphericalAberrationCorrector struct {
	/*
	 * whether a specific instrument was used during data acquisition
	 */
	Used basetypes.Bool `json:"used"`
	/*
	 * Details of a given specialist instrument
	 */
	InstrumentType basetypes.String `json:"instrument_type"`
}

/*
 * Special device used to correct instrument inherent chromatic aberration.
 */
type ChromaticAberrationCorrector struct {
	/*
	 * whether a specific instrument was used during data acquisition
	 */
	Used basetypes.Bool `json:"used"`
	/*
	 * Details of a given specialist instrument
	 */
	InstrumentType basetypes.String `json:"instrument_type"`
}

/*
 * Instrument values, mostly constant across a data collection.
 */
type Instrument struct {
	/*
	 * Name/Type of the Microscope
	 */
	Microscope basetypes.String `json:"microscope"`
	/*
	 * Mode of illumination used during data collection
	 */
	Illumination basetypes.String `json:"illumination"`
	/*
	 * Mode of imaging used during data collection
	 */
	Imaging basetypes.String `json:"imaging"`
	/*
	 * Type of electron source used in the microscope, such as FEG
	 */
	ElectronSource basetypes.String `json:"electron_source"`
	/*
	 * Voltage used for the electron acceleration, in kV
	 */
	AccelerationVoltage basetypes.Float64 `json:"acceleration_voltage"`
	/*
	 * C2 aperture size used in data acquisition, in µm
	 */
	C2Aperture basetypes.Float64 `json:"c2_aperture"`
	/*
	 * Spherical aberration of the instrument, in mm
	 */
	Cs basetypes.Float64 `json:"cs"`
}

/*
 * Unifying class to describe the full sample.
 */
type Sample struct {
	/*
	 * Description of the overall molecule
	 */
	OverallMolecule OverallMolecule `json:"overall_molecule"`
	/*
	 * List of molecule associated with the sample
	 */
	Molecule []Molecule `json:"molecule"`
	/*
	 * List of ligands associated with the sample
	 */
	Ligands []Ligand `json:"ligands"`
	/*
	 * Description of the specimen
	 */
	Specimen Specimen `json:"specimen"`
	/*
	 * Description of the grid used
	 */
	Grid Grid `json:"grid"`
}

/*
 * Description of the overall molecule
 */
type OverallMolecule struct {
	/*
	 * Description of the overall molecular type, i.e., a complex
	 */
	MolecularType basetypes.String `json:"molecular_type"`
	/*
	 * Name of the full sample
	 */
	NameSample basetypes.String `json:"name_sample"`
	/*
	 * Where the sample was taken from, i.e., natural host, recombinantly expressed, etc.
	 */
	Source basetypes.String `json:"source"`
	/*
	 * Molecular weight in Da
	 */
	MolecularWeight basetypes.Float64 `json:"molecular_weight"`
	/*
	 * What type of higher order structure your sample forms - if any.
	 */
	Assembly basetypes.String `json:"assembly"`
}

/*
 * More detailed information about individual molecules.
 */
type Molecule struct {
	/*
	 * Name of an individual molecule (often protein) in the sample
	 */
	NameMol basetypes.String `json:"name_mol"`
	/*
	 * Description of the overall molecular type, i.e., a complex
	 */
	MolecularType basetypes.String `json:"molecular_type"`
	/*
	 * Class of the molecule
	 */
	MolecularClass basetypes.String `json:"molecular_class"`
	/*
	 * Full sequence of the sample as in the data, i.e., cleaved tags should also be removed from sequence here
	 */
	Sequence basetypes.String `json:"sequence"`
	/*
	 * Scientific name of the natural host organism
	 */
	NaturalSource basetypes.String `json:"natural_source"`
	/*
	 * Taxonomy ID of the natural source organism
	 */
	TaxonomyIdSource basetypes.String `json:"taxonomy_id_source"`
	/*
	 * Scientific name of the organism used to produce the molecule of basetypes.Interest
	 */
	ExpressionSystem basetypes.String `json:"expression_system"`
	/*
	 * Taxonomy ID of the expression system organism
	 */
	TaxonomyIdExpression basetypes.String `json:"taxonomy_id_expression"`
	/*
	 * Name of the gene of basetypes.Interest
	 */
	GeneName basetypes.String `json:"gene_name"`
}

/*
 * Information on ligands if present.
 */
type Ligand struct {
	/*
	 * Whether the model contains any ligands
	 */
	Present basetypes.Bool `json:"present"`
	/*
	 * Provide a valid SMILES basetypes.String of your ligand
	 */
	Smiles basetypes.String `json:"smiles"`
	/*
	 * Link to a reference of your ligand, i.e., CCD, PubChem, etc.
	 */
	Reference basetypes.String `json:"reference"`
}

/*
 * Description of specimen handling.
 */
type Specimen struct {
	/*
	 * Name/composition of the (chemical) sample buffer during grid preparation
	 */
	Buffer basetypes.String `json:"buffer"`
	/*
	 * Concentration of the (supra)molecule in the sample, in mg/ml
	 */
	Concentration basetypes.Float64 `json:"concentration"`
	/*
	 * pH of the sample buffer
	 */
	Ph basetypes.Float64 `json:"ph"`
	/*
	 * Whether the sample was vitrified
	 */
	Vitrification basetypes.Bool `json:"vitrification"`
	/*
	 * Which cryogen was used for vitrification
	 */
	VitrificationCryogen basetypes.String `json:"vitrification_cryogen"`
	/*
	 * Environmental humidity just before vitrification, in %
	 */
	Humidity basetypes.Float64 `json:"humidity"`
	/*
	 * Environmental temperature just before vitrification, in K
	 */
	Temperature basetypes.Float64 `json:"temperature"`
	/*
	 * Whether the sample was stained
	 */
	Staining basetypes.Bool `json:"staining"`
	/*
	 * Whether the sample was embedded
	 */
	Embedding basetypes.Bool `json:"embedding"`
	/*
	 * Whether the sample was shadowed
	 */
	Shadowing basetypes.Bool `json:"shadowing"`
}

/*
 * Details on the grid used in the experiment.
 */
type Grid struct {
	/*
	 * Grid manufacturer
	 */
	Manufacturer basetypes.String `json:"manufacturer"`
	/*
	 * Material out of which the grid is made
	 */
	Material basetypes.String `json:"material"`
	/*
	 * Grid mesh in lines per inch
	 */
	Mesh basetypes.Float64 `json:"mesh"`
	/*
	 * Whether a support film was used
	 */
	FilmSupport basetypes.Bool `json:"film_support"`
	/*
	 * Type of material the support film is made of
	 */
	FilmMaterial basetypes.String `json:"film_material"`
	/*
	 * Topology of the support film
	 */
	FilmTopology basetypes.String `json:"film_topology"`
	/*
	 * Thickness of the support film
	 */
	FilmThickness basetypes.String `json:"film_thickness"`
	/*
	 * Type of pretreatment of the grid, i.e., glow discharge
	 */
	PretreatmentType basetypes.String `json:"pretreatment_type"`
	/*
	 * Length of time of the pretreatment in s
	 */
	PretreatmentTime basetypes.Float64 `json:"pretreatment_time"`
	/*
	 * Pressure of the chamber during pretreatment, in Pa
	 */
	PretreatmentPressure basetypes.Float64 `json:"pretreatment_pressure"`
	/*
	 * Atmospheric conditions in the chamber during pretreatment, i.e., addition of specific gases, etc.
	 */
	PretreatmentAtmosphere basetypes.String `json:"pretreatment_atmosphere"`
}

/*
 * The min, max and increment of the tilt angle in a tomography session. Unit is degree.
 */
type TiltAngle struct {
	/*
	 * parent types
	 */
	Series
	/*
	 * Increment between elements of a series
	 */
	Increment basetypes.Float64 `json:"increment"`
	/*
	 * Minimal value of a given dataset property
	 */
	Minimal basetypes.Float64 `json:"minimal"`
	/*
	 * Maximal value of a given dataset property
	 */
	Maximal basetypes.Float64 `json:"maximal"`
}


type AcquisitionTomo struct {
	/*
	 * parent types
	 */
	//Acquisition
	/*
	 * The tilt axis angle of a tomography series
	 */
	TiltAxisAngle basetypes.Float64 `json:"tilt_axis_angle"`
	/*
	 * The min, max and increment of the tilt angle in a tomography session. Unit is degree.
	 */
	TiltAngle TiltAngle `json:"tilt_angle"`
	/*
	 * Target defocus set, min and max values in µm.
	 */
	NominalDefocus Range `json:"nominal_defocus"`
	/*
	 * Machine estimated defocus, min and max values in µm. Has a tendency to be off.
	 */
	CalibratedDefocus Range `json:"calibrated_defocus"`
	/*
	 * Magnification level as indicated by the instrument, no unit
	 */
	NominalMagnification basetypes.Int `json:"nominal_magnification"`
	/*
	 * Calculated magnification, no unit
	 */
	CalibratedMagnification basetypes.Int `json:"calibrated_magnification"`
	/*
	 * Speciman holder model
	 */
	Holder basetypes.String `json:"holder"`
	/*
	 * Type of cryogen used in the holder - if the holder is cooled seperately
	 */
	HolderCryogen basetypes.String `json:"holder_cryogen"`
	/*
	 * Temperature during data collection, in K with min and max values.
	 */
	Temperature Range `json:"temperature"`
	/*
	 * Software used for instrument control,
	 */
	MicroscopeSoftware basetypes.String `json:"microscope_software"`
	/*
	 * Make and model of the detector used
	 */
	Detector basetypes.String `json:"detector"`
	/*
	 * Operating mode of the detector
	 */
	DetectorMode basetypes.String `json:"detector_mode"`
	/*
	 * Average dose per image/movie/tilt - given in electrons per square Angstrom
	 */
	DosePerMovie basetypes.Float64 `json:"dose_per_movie"`
	/*
	 * Wether an energy filter was used and its specifics.
	 */
	EnergyFilter EnergyFilter `json:"energy_filter"`
	/*
	 * The size of the image in pixels, height and width given.
	 */
	ImageSize ImageSize `json:"image_size"`
	/*
	 * Time and date of the data acquisition
	 */
	DateTime basetypes.String `json:"date_time"`
	/*
	 * Time of data acquisition per movie/tilt - in s
	 */
	ExposureTime basetypes.Float64 `json:"exposure_time"`
	/*
	 * Cryogen used in cooling the instrument and sample, usually nitrogen
	 */
	Cryogen basetypes.String `json:"cryogen"`
	/*
	 * Number of frames that on average constitute a full movie, can be a bit hard to define for some detectors
	 */
	FramesPerMovie basetypes.Int `json:"frames_per_movie"`
	/*
	 * Number of grids imaged for this project - here with qualifier during this data acquisition
	 */
	GridsImaged basetypes.Int `json:"grids_imaged"`
	/*
	 * Number of images generated total for this data collection - might need a qualifier for tilt series to determine whether full series or individual tilts are counted
	 */
	ImagesGenerated basetypes.Int `json:"images_generated"`
	/*
	 * Level of binning on the images applied during data collection
	 */
	BinningCamera basetypes.Float64 `json:"binning_camera"`
	/*
	 * Pixel size, in Angstrom
	 */
	PixelSize basetypes.Float64 `json:"pixel_size"`
	/*
	 * basetypes.String type of special optics, such as a phaseplate
	 */
	SpecialistOptics SpecialistOptics `json:"specialist_optics"`
	/*
	 * Movement of the beam above the sample for data collection purposes that does not require movement of the stage. Given in mrad.
	 */
	Beamshift BoundingBox2D `json:"beamshift"`
	/*
	 * Another way to move the beam above the sample for data collection purposes that does not require movement of the stage. Given in mrad.
	 */
	Beamtilt BoundingBox2D `json:"beamtilt"`
	/*
	 * Movement of the Beam below the image in order to shift the image on the detector. Given in µm.
	 */
	Imageshift BoundingBox2D `json:"imageshift"`
	/*
	 * Number of Beamtilt groups present in this dataset - for optimized processing split dataset basetypes.Into groups of same tilt angle. Despite its name Beamshift is often used to achive this result.
	 */
	Beamtiltgroups basetypes.Int `json:"beamtiltgroups"`
	/*
	 * Whether and how you have to flip or rotate the gainref in order to align with your acquired images
	 */
	GainrefFlipRotate basetypes.String `json:"gainref_flip_rotate"`
}

/*
 * Overarching category for authors and grants
 */
type Organizational struct {
	/*
	 * List of grants associated with the project
	 */
	Grants []Grant `json:"grants"`
	/*
	 * List of authors associated with the project
	 */
	Authors []Author `json:"authors"`
	/*
	 * funding organization/person.
	 */
	Funder []Funder `json:"funder"`
}

/*
 * personal information
 */
type Person struct {
	/*
	 * name
	 */
	Name basetypes.String `json:"name"`
	/*
	 * first name
	 */
	FirstName basetypes.String `json:"first_name"`
	/*
	 * work status
	 */
	WorkStatus basetypes.Bool `json:"work_status"`
	/*
	 * email
	 */
	Email basetypes.String `json:"email"`
	/*
	 * work phone
	 */
	WorkPhone basetypes.String `json:"work_phone"`
}

/*
 * Details on the person performing the experiment.
 */
type Author struct {
	/*
	 * parent types
	 */
	Person
	/*
	 * ORCID of the author, a type of unique identifier
	 */
	Orcid basetypes.String `json:"orcid"`
	/*
	 * Country of the institution
	 */
	Country basetypes.String `json:"country"`
	/*
	 * Role of the author, for example principal investigator
	 */
	Role basetypes.String `json:"role"`
	/*
	 * Name of the organization
	 */
	NameOrg basetypes.String `json:"name_org"`
	/*
	 * Type of organization, academic, commercial, governmental, etc.
	 */
	TypeOrg basetypes.String `json:"type_org"`
	/*
	 * name
	 */
	Name basetypes.String `json:"name"`
	/*
	 * first name
	 */
	FirstName basetypes.String `json:"first_name"`
	/*
	 * work status
	 */
	WorkStatus basetypes.Bool `json:"work_status"`
	/*
	 * email
	 */
	Email basetypes.String `json:"email"`
	/*
	 * work phone
	 */
	WorkPhone basetypes.String `json:"work_phone"`
}

/*
 * Grant
 */
type Grant struct {
	/*
	 * name of the grant
	 */
	GrantName basetypes.String `json:"grant_name"`
	/*
	 * start date
	 */
	StartDate basetypes.String `json:"start_date"`
	/*
	 * end date
	 */
	EndDate basetypes.String `json:"end_date"`
	/*
	 * budget
	 */
	Budget basetypes.Float64 `json:"budget"`
	/*
	 * project id
	 */
	ProjectId basetypes.String `json:"project_id"`
	/*
	 * Country of the institution
	 */
	Country basetypes.String `json:"country"`
}

/*
 * Description of the project funding
 */
type Funder struct {
	/*
	 * funding organization/person.
	 */
	FunderName basetypes.String `json:"funder_name"`
	/*
	 * Type of organization, academic, commercial, governmental, etc.
	 */
	TypeOrg basetypes.String `json:"type_org"`
	/*
	 * Country of the institution
	 */
	Country basetypes.String `json:"country"`
}

/*
 * OSC-EM Metadata for a dataset
 */
type EMDatasetBase struct {
	//Acquisition basetypes.String `json:"acquisition"`
	Instrument basetypes.String `json:"instrument"`
	Sample basetypes.String `json:"sample"`
	Organizational basetypes.String `json:"organizational"`
}

/*
 * cryo electron tomography dataset, with focus on a single protein (potentially subtomogram averaging)
 */
type EMDatasetTomo struct {
	/*
	 * parent types
	 */
	EMDatasetBase
	/*
	 * Describe the data acquisition parameters
	 */
	//Acquisition //AcquisitionTomo `json:"acquisition"`
	/*
	 * Describe the instrument used to acquire the data
	 */
	Instrument Instrument `json:"instrument"`
	/*
	 * Sample information
	 */
	Sample Sample `json:"sample"`
	/*
	 * Information on authors and grants
	 */
	Organizational Organizational `json:"organizational"`
}


