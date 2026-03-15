import React, { useState, useCallback } from 'react';
import { Calendar, momentLocalizer, EventProps } from 'react-big-calendar';
//import withDragAndDrop from 'react-big-calendar/lib/addons/dragAndDrop';
import 'react-big-calendar/lib/css/react-big-calendar.css';
import 'react-big-calendar/lib/addons/dragAndDrop/styles.css';
import { format, parse, startOfWeek, getDay} from 'date-fns';
import { ru } from 'date-fns/locale/ru';
import { dateFnsLocalizer } from 'react-big-calendar';
import Modal from 'react-modal';
import { CalendarEvent } from './types';
import './styles.css';

// Локализация на русский (date-fns)
//const locales = {
//  'ru': require('date-fns/locale/ru'),
//};

const locales = {ru};
const localizer = dateFnsLocalizer({
  format,
  parse,
  startOfWeek,
  getDay,
  locales,
});

// Обёртка для Drag-and-Drop
//const DragAndDropCalendar = withDragAndDrop(Calendar);
const DragAndDropCalendar = Calendar;

Modal.setAppElement('#root');

const CustomCalendar: React.FC = () => {
  const [events, setEvents] = useState<CalendarEvent[]>([
    {
      id: '1',
      title: 'Встреча с заказчиком',
      start: new Date(2026, 2, 16, 10, 0),
      end: new Date(2026, 2, 16, 11, 30),
      author: 'Иван Петров',
      color: '#3174ad',
    },
    {
      id: '2',
      title: 'Разработка',
      start: new Date(2026, 2, 17, 14, 0),
      end: new Date(2026, 2, 17, 16, 0),
      author: 'Мария Иванова',
      color: '#e57373',
    },
  ]);

  const [modalIsOpen, setModalIsOpen] = useState(false);
  const [selectedEvent, setSelectedEvent] = useState<CalendarEvent | null>(null);
  const [newEvent, setNewEvent] = useState<Partial<CalendarEvent>>({
    title: '',
    author: '',
    color: '#3174ad',
  });

  // Открыть модалку для создания события
  const handleSelectSlot = useCallback(({ start, end }: { start: Date; end: Date }) => {
    setNewEvent({
      title: '',
      author: '',
      color: '#3174ad',
      start,
      end,
    });
    setSelectedEvent(null);
    setModalIsOpen(true);
  }, []);

  // Открыть модалку для редактирования/удаления
  const handleSelectEvent = useCallback((event: CalendarEvent) => {
    setSelectedEvent(event);
    setNewEvent(event);
    setModalIsOpen(true);
  }, []);

  // Перемещение события (drag & drop)
  const handleEventDrop = useCallback(({ event, start, end }: any) => {
    setEvents((prev) =>
      prev.map((e) => (e.id === event.id ? { ...e, start, end } : e))
    );
  }, []);

  // Изменение размера события
  const handleEventResize = useCallback(({ event, start, end }: any) => {
    setEvents((prev) =>
      prev.map((e) => (e.id === event.id ? { ...e, start, end } : e))
    );
  }, []);

  // Сохранить событие (новое или изменённое)
  const handleSaveEvent = () => {
    if (!newEvent.title || !newEvent.start || !newEvent.end) {
      alert('Заполните все поля');
      return;
    }

    if (selectedEvent) {
      // Редактирование существующего
      setEvents((prev) =>
        prev.map((e) =>
          e.id === selectedEvent.id ? ({ ...selectedEvent, ...newEvent } as CalendarEvent) : e
        )
      );
    } else {
      // Создание нового
      const event: CalendarEvent = {
        id: Math.random().toString(36).substr(2, 9),
        title: newEvent.title!,
        start: newEvent.start!,
        end: newEvent.end!,
        author: newEvent.author!,
        color: newEvent.color!,
      };
      setEvents([...events, event]);
    }
    setModalIsOpen(false);
  };

  // Удалить событие
  const handleDeleteEvent = () => {
    if (selectedEvent) {
      setEvents(events.filter((e) => e.id !== selectedEvent.id));
      setModalIsOpen(false);
    }
  };

  // Пропсы для стилизации событий (цвет)
  const eventPropGetter = useCallback(
    (event: CalendarEvent) => ({
      style: {
        backgroundColor: event.color,
        borderRadius: '4px',
        opacity: 0.8,
        color: 'white',
        border: '0',
        display: 'block',
      },
    }),
    []
  );
// Кастомный компонент события для отображения автора
  const EventComponent: React.FC<EventProps<CalendarEvent>> = ({ event }) => (
    <div>
      <strong>{event.title}</strong>
      <br />
      <small>Автор: {event.author}</small>
    </div>
  );

  return (
    <div className="app">
      <h1>Календарь задач</h1>
      <DragAndDropCalendar
        localizer={localizer}
        events={events}
        startAccessor="start"
        endAccessor="end"
        style={{ height: "80vh", margin: '40px', backgroundColor:"#fff" }}
        selectable
        onSelectSlot={handleSelectSlot}
        onSelectEvent={handleSelectEvent}
//        onEventDrop={handleEventDrop}
//        onEventResize={handleEventResize}
        eventPropGetter={eventPropGetter}
        components={{
          event: EventComponent,
        }}
        defaultView="week"
        views={['month', 'week', 'day', 'agenda']}
        messages={{
          week: 'Неделя',
          day: 'День',
          month: 'Месяц',
          agenda: 'Повестка',
          today: 'Сегодня',
          previous: 'Назад',
          next: 'Вперёд',
        }}
      />

      {/* Модальное окно для создания/редактирования */}
      <Modal
        isOpen={modalIsOpen}
        onRequestClose={() => setModalIsOpen(false)}
        contentLabel="Событие"
        className="customModal"
        overlayClassName="customOverlay2"
      >
        <h2>{selectedEvent ? 'Редактировать событие' : 'Новое событие'}</h2>
        <form
          onSubmit={(e) => {
            e.preventDefault();
            handleSaveEvent();
          }}
        >
          
          <div className='row mb-2'>
            <div className='col-3'>Название:</div>
            <div className='col-9'>
            <input
              type="text"
              value={newEvent.title || ''}
              onChange={(e) => setNewEvent({ ...newEvent, title: e.target.value })}
              required
            /></div>
          </div>
          <div className='row mb-2'>
            <div className='col-3'>Автор:</div>
            <div className='col-9'>
            <input
              type="text"
              value={newEvent.author || ''}
              onChange={(e) => setNewEvent({ ...newEvent, author: e.target.value })}
              required
            /></div>
          </div>
          <div className='row mb-2'>
            <div className='col-3'>Цвет:</div>
            <div className='col-9'>
            <input
              type="color"
              value={newEvent.color || '#3174ad'}
              onChange={(e) => setNewEvent({ ...newEvent, color: e.target.value })}
            /></div>
          </div>
          <div className='row mb-2'>
            <div className='col-3'>Начало:</div>
            <div className='col-9'>
            <input
              type="datetime-local"
              value={newEvent.start ? format(newEvent.start, "yyyy-MM-dd'T'HH:mm") : ''}
              onChange={(e) =>
                setNewEvent({ ...newEvent, start: e.target.value ? new Date(e.target.value) : undefined })
              }
              required
            /></div>
          </div>
          <div className='row mb-2'>
            <div className='col-3'>Конец:</div>
            <div className='col-9'>
            <input
              type="datetime-local"
              value={newEvent.end ? format(newEvent.end, "yyyy-MM-dd'T'HH:mm") : ''}
              onChange={(e) =>
                setNewEvent({ ...newEvent, end: e.target.value ? new Date(e.target.value) : undefined })
              }
              required
            /></div>
          </div>
          <div className="modal-actions">
            <button type="submit">Сохранить</button>
            {selectedEvent && (
              <button type="button" onClick={handleDeleteEvent} className="delete">
                Удалить
              </button>
            )}
            <button type="button" onClick={() => setModalIsOpen(false)}>
              Отмена
            </button>
          </div>
        </form>
      </Modal>
    </div>
  );
};

export default CustomCalendar;
