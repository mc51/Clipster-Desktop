#ifndef GOEY_THUNKS_H
#define GOEY_THUNKS_H

#include <stdbool.h>
#include <stddef.h>

extern void widgetClose( void *widget );
extern bool widgetCanFocus( void *widget );
extern void widgetGrabFocus( void *widget );
extern bool widgetIsFocus( void *widget );
extern void widgetSendKey( void *widget, unsigned key, bool release );
extern void widgetNaturalSize( void *widget, int *width, int *height );
extern int widgetMinHeight( void *widget );
extern int widgetMinHeightForWidth( void *widget, int width );
extern int widgetNaturalHeight( void *widget );
extern int widgetNaturalHeightForWidth( void *widget, int width );
extern int widgetMinWidth( void *widget );
extern int widgetNaturalWidth( void *widget );
extern int widgetNaturalWidthForHeight( void *widget, int height );
extern void widgetSetBounds( void *widget, int x, int y, int width,
                             int height );
extern bool widgetSensitive( void *widget );
extern bool widgetCanDefault( void *widget );
extern void widgetSetSizeRequest( void *widget, int width, int height );

typedef struct {
    int width;
    int height;
} windowsize_t;

extern void *mountWindow( char const *text );
extern windowsize_t windowSize( void *window );
extern void *windowScrolledWindow( void *window );
extern void *windowLayout( void *window );
extern void windowSetLayoutSize( void *window, unsigned width,
                                 unsigned height );
extern char const *windowTitle( void *window );
extern void windowSetTitle( void *window, char const *text );
extern unsigned windowVScrollbarWidth( void *window );
extern unsigned windowHScrollbarHeight( void *window );
extern void windowShowScrollbars( void *window, bool horz, bool vert );
extern void windowShow( void *window );
extern void windowSetDefaultSize( void *window, int width, int height );
extern void windowSetIcon( void *window, unsigned char const *data, int width,
                           int height, int rowStride );
extern void *windowScreenshot( void *window, void **data, size_t *datLen,
                               bool *haveAlpha, int *width, int *height,
                               unsigned *stride );

extern unsigned dialogRun( void *dialog );
extern void dialogAddFilter( void *dialog, char const *name,
                             char const *pattern );
extern unsigned dialogResponseAccept( void );
extern unsigned dialogResponseCancel( void );

extern void *mountMessageDialog( void *window, char const *title, unsigned icon,
                                 char const *text );
extern unsigned messageDialogWithError( void );
extern unsigned messageDialogWithWarn( void );
extern unsigned messageDialogWithInfo( void );

extern void *mountOpenDialog( void *dialog, char const *title,
                              char const *filename );
extern void *mountSaveDialog( void *dialog, char const *title,
                              char const *filename );
extern char const *dialogGetFilename( void *dialog );

extern void *mountButton( void *container, char const *text, bool disabled,
                          bool def, bool onclick, bool onfocus, bool onblur );
extern void buttonUpdate( void *button, char const *text, bool disabled,
                          bool def, bool onclick, bool onfocus, bool onblur );
extern void buttonClick( void *button );
extern char const *buttonText( void *button );

extern void *mountLabel( void *container, char const *text );
extern void labelUpdate( void *label, char const *text );
extern char const *labelText( void *label );

extern void *mountCheckbox( void *container, bool value, char const *text,
                            bool disabled, bool onchange, bool onfocus,
                            bool onblur );
extern void checkboxUpdate( void *button, bool value, char const *text,
                            bool disabled, bool onchange, bool onfocus,
                            bool onblur );
extern void checkboxClick( void *button );
extern bool checkboxValue( void *button );
extern char const *checkboxText( void *button );

extern void *mountTextbox( void *container, char const *text,
                           char const *placeholder, bool disabled,
                           bool password, bool readonly, bool onchange,
                           bool onfocus, bool onblur, bool onenterkey );
extern void textboxUpdate( void *widget, char const *text,
                           char const *placeholder, bool disabled,
                           bool password, bool readonly, bool onchange,
                           bool onfocus, bool onblur, bool onenterkey );
extern char const *textboxText( void *widget );
extern char const *textboxPlaceholder( void *widget );
extern bool textboxPassword( void *widget );
extern bool textboxReadOnly( void *widget );

extern void *mountTextarea( void *container, char const *text, bool disabled,
                            bool readonly, bool onchange, bool onfocus,
                            bool onblur );
extern void textareaUpdate( void *widget, char const *text, bool disabled,
                            bool readonly, bool onchange, bool onfocus,
                            bool onblur );
extern char const *textareaText( void *widget );
extern char const *textareaPlaceholder( void *widget );
extern bool textareaReadOnly( void *widget );
extern void *textareaTextview( void *widget );

extern void *mountParagraph( void *container, char const *text, char align );
extern void paragraphUpdate( void *widget, char const *text, char align );
extern char const *paragraphText( void *widget );
extern char paragraphAlign( void *widget );
extern void paragraphSetText( void *widget, char const *text );

extern void *mountProgressbar( void *contianer, double value );
extern void progressbarUpdate( void *widget, double value );
extern double progressbarValue( void *widget );

extern void *mountSlider( void *container, double value, bool disabled,
                          double min, double max, bool onchange, bool onfocus,
                          bool onblur );
extern void sliderUpdate( void *widget, double value, bool disabled, double min,
                          double max, bool onchange, bool onfocus,
                          bool onblur );
extern double sliderValue( void *widget );

extern void *mountHR( void *parent );

extern void *mountImage( void *parent, unsigned char const *data, int width,
                         int height, int rowStride );
extern void *imageUpdate( void *widget, unsigned char const *data, int width,
                          int height, int rowStride );
extern unsigned imageColorSpace( void *widget );
extern bool imageHasAlpha( void *widget );
extern unsigned imageImageWidth( void *widget );
extern unsigned imageImageHeight( void *widget );
extern unsigned imageImageStride( void *widget );
extern char const *imageImageData( void *widget, size_t *length );

extern void *mountDecoration( void *parent, unsigned fill, unsigned stroke,
                              int radius );
extern void decorationUpdate( void *widget, unsigned fill, unsigned stroke,
                              int radius );
extern unsigned decorationFill( void *widget );
extern unsigned decorationStroke( void *widget );
extern int decorationRadius( void *widget );
extern void decorationSetRadius( void *widget, int radius );

extern void *mountCombobox( void *parent, char const *items, int value,
                            bool unset, bool disabled, bool onchange,
                            bool onfocus, bool onblur );
extern void comboboxUpdate( void *widget, char const *items, int value,
                            bool unset, bool disabled, bool onchange,
                            bool onfocus, bool onblur );
extern void *comboboxChild( void *widget );
extern unsigned comboboxItemCount( void *widget );
extern char const *comboboxItem( void *widget, unsigned index );
extern int comboboxValue( void *widget );

extern void *mountTabs( void *parent, int value, char *tabs, bool onchange );
extern void tabsUpdate( void *widget, int value, char *tabs, bool onchange );
extern void *tabsGetTabParent( void *widget, int value );
extern int tabsItemCount( void *widget );
extern char const *tabsItemCaption( void *widget, int index );

extern void *mountIntInput( void *parent, long value, char const *placeholder,
                            bool disabled, long min, long max, bool onchange,
                            bool onfocus, bool onblur, bool onenterkey );
extern void intinputUpdate( void *widget, long value, char const *placeholder,
                            bool disabled, long min, long max, bool onchange,
                            bool onfocus, bool onblur, bool onenterkey );
extern long intinputValue( void *widget );
extern double intinputMin( void *widget );
extern double intinputMax( void *widget );

extern void *mountDateInput( void *parent, int year, unsigned month,
                             unsigned day, bool disabled, bool onchange,
                             bool onfocus, bool onblur );
extern void dateInputUpdate( void *widget, int year, unsigned month,
                             unsigned day, bool disabled, bool onchange,
                             bool onfocus, bool onblur );
extern int dateInputYear( void *widget );
extern unsigned dateInputMonth( void *widget );
extern unsigned dateInputDay( void *widget );

#endif