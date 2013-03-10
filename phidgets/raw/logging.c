#include "logging.h"

int _log(CPhidgetLog_level l, const char * id, const char * message) {
  return CPhidget_log(l, id, message);
}
